package gogen

import (
	_ "embed"
	"fmt"
	"github.com/zhangxiaohan228/zero-admin-goctl/api/parser/g4/gen/api"
	"github.com/zhangxiaohan228/zero-admin-goctl/api/spec"
	"github.com/zhangxiaohan228/zero-admin-goctl/config"
	"github.com/zhangxiaohan228/zero-admin-goctl/util/format"
	"github.com/zhangxiaohan228/zero-admin-goctl/util/pathx"
	"github.com/zhangxiaohan228/zero-admin-goctl/vars"
	"path"
	"strconv"
	"strings"
)

//go:embed logic.tpl
var logicTemplate string

var (
	genLogicFunc = map[string]string{
		"/list": "FindList",
		"/:id":  "FindOne",
	}
	genLogicFuncParam = map[string]string{
		"/list": "l.ctx,\"*\",nil",
		"/:id":  "l.ctx,req.Id",
	}
	// 表示响应值是否需要经过结构体赋值
	genLogicFuncResult = map[string]bool{
		"FindList":         true,
		"FindListWithPage": true,
		"FindOne":          false,
	}
)

func genLogic(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	groupModel := getModel(api)

	for _, g := range api.Service.Groups {
		serviceGroup := g.GetAnnotation("group")
		model := groupModel[serviceGroup]
		for _, r := range g.Routes {
			err := genLogicByRoute(dir, rootPkg, cfg, g, r, newLogicHandlerByMethod(model, api.Service.Groups))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByRoute(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route, logicHandle map[string]spec.LogicHandle) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	imports := genLogicImports(route, rootPkg)
	if logicHandle[logic].Import != "" {
		imports += logicHandle[logic].Import
	}

	var responseString string
	var returnString string
	var requestString string
	if len(route.ResponseTypeName()) > 0 {
		resp := responseGoTypeName(route, typesPacket)
		responseString = "(resp " + resp + ", err error)"
		returnString = "return"
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestTypeName()) > 0 {
		requestString = "req *" + requestGoTypeName(route, typesPacket)
	}

	subDir := getLogicFolderPath(group, route)
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          subDir,
		filename:        goFile + ".go",
		templateName:    "logicTemplate",
		category:        category,
		templateFile:    logicTemplateFile,
		builtinTemplate: logicTemplate,
		data: map[string]any{
			"pkgName":      subDir[strings.LastIndex(subDir, "/")+1:],
			"imports":      imports,
			"logic":        strings.Title(logic),
			"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
			"responseType": responseString,
			"returnString": returnString,
			"request":      requestString,
			"hasDoc":       len(route.JoinedDoc()) > 0,
			"doc":          getDoc(route.JoinedDoc()),
			"logicHandle":  genLogicHandle(logicHandle[logic], logic, responseString),
		},
	})
}

func getLogicFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return logicDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDir, folder)
}

func genLogicImports(route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, `"context"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)))
	if shallImportTypesPackage(route) {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}

func onlyPrimitiveTypes(val string) bool {
	fields := strings.FieldsFunc(val, func(r rune) bool {
		return r == '[' || r == ']' || r == ' '
	})

	for _, field := range fields {
		if field == "map" {
			continue
		}
		// ignore array dimension number, like [5]int
		if _, err := strconv.Atoi(field); err == nil {
			continue
		}
		if !api.IsBasicType(field) {
			return false
		}
	}

	return true
}

func shallImportTypesPackage(route spec.Route) bool {
	if len(route.RequestTypeName()) > 0 {
		return true
	}

	respTypeName := route.ResponseTypeName()
	if len(respTypeName) == 0 {
		return false
	}

	if onlyPrimitiveTypes(respTypeName) {
		return false
	}

	return true
}

func genLogicHandle(logicHandle spec.LogicHandle, logic, responseString string) string {

	if logicHandle.Model == "" {
		return ""
	}

	var handles []string

	modelStruct := strings.Title(logicHandle.Model) + "Model"

	funcResult := logicHandle.Result
	if strings.Contains(logicHandle.Result, "_") {
		funcResult += " = "
	} else {
		funcResult += ":="
	}
	if logicHandle.ReqIsModel {
		handles = append(handles, fmt.Sprintf("var data *%s.%s", logicHandle.Model, strings.Title(logicHandle.Model)))
		handles = append(handles, "copier.Copy(&data, &req)"+"\n")
		handles = append(handles, fmt.Sprintf("%s l.svcCtx.Models.%s.%s(l.ctx,data)", funcResult, modelStruct, logicHandle.Function))
	} else {
		handles = append(handles, fmt.Sprintf("%s l.svcCtx.Models.%s.%s(%s)", funcResult, modelStruct, logicHandle.Function, logicHandle.Params))
	}
	handles = append(handles, fmt.Sprintf("if err != nil{\n\t\treturn nil,err\n\t}"))

	if logicHandle.ResultIsModel {
		parts := strings.SplitN(responseString, "*", 2)
		typeName := strings.SplitN(parts[1], ", ", 2)[0]
		handles = append(handles, fmt.Sprintf("resp = &%s{}", typeName))

		if genLogicFuncResult[logicHandle.Function] {
			handles = append(handles, fmt.Sprintf("var data []*types.%s", logic[:strings.Index(logic, "Logic")]))
			handles = append(handles, "copier.Copy(&data, &result)"+"\n")
			handles = append(handles, "resp.List = data")
		} else {
			handles = append(handles, "copier.Copy(&resp, &result)"+"\n")
		}
	}

	return strings.Join(handles, "\n\t")
}

func newLogicHandlerByMethod(model string, apiGroups []spec.Group) map[string]spec.LogicHandle {

	if model == "" {
		return nil
	}

	logicModelInterface := make(map[string]spec.LogicHandle, 0)

	for _, g := range apiGroups {
		for _, r := range g.Routes {
			if r.Path != "/" {
				// 判断是否为需要生成的路由
				if _, ok := genLogicFunc[r.Path]; !ok {
					continue
				}
			}
			logicName := getLogicName(r)
			if r.Path == "/" {
				switch r.Method {
				case "get":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:      model,
						Function:   "FindOne",
						Params:     "l.ctx,req.Id",
						Result:     "result,err",
						ReqIsModel: false,
						Import:     "\n\"github.com/jinzhu/copier\"\n",
					}
				case "post":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:      model,
						Function:   "Insert",
						Params:     "l.ctx",
						Result:     "_,err",
						ReqIsModel: true,
						Import:     fmt.Sprintf("\n\"github.com/jinzhu/copier\"\n\t \"zero-admin/server/model/%s\"", model),
					}
				case "put":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:      model,
						Function:   "Update",
						Params:     "l.ctx",
						Result:     "_,err",
						ReqIsModel: true,
						Import:     fmt.Sprintf("\n\"github.com/jinzhu/copier\"\n\t \"zero-admin/server/model/%s\"", model),
					}
				case "delete":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:      model,
						Function:   "Delete",
						Params:     "l.ctx,req.Id",
						Result:     "_,err",
						ReqIsModel: false,
					}
				}
			} else {
				logicModelInterface[logicName] = spec.LogicHandle{
					Model:         model,
					Function:      genLogicFunc[r.Path],
					Params:        genLogicFuncParam[r.Path],
					Result:        "result,err",
					ResultIsModel: true,
					Import:        "\n\"github.com/jinzhu/copier\"\n",
				}
			}
		}
	}

	return logicModelInterface
}
