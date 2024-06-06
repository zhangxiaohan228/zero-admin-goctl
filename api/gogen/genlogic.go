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

func genLogic(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {

	logicHandle := newLogicHandlerByMethod(getModel(api), api.Service.Groups)
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicByRoute(dir, rootPkg, cfg, g, r, logicHandle)
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
			"logicHandler": genLogicHandle(logicHandle[logic]),
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

func genLogicHandle(logicHandle spec.LogicHandle) string {

	var handles []string

	modelStruct := strings.Title(logicHandle.Model) + "Model"
	if logicHandle.IsModel {
		handles = append(handles, fmt.Sprintf("var data %s.%s", logicHandle.Model, modelStruct))
		handles = append(handles, "copier.Copy(&data, &req)"+"\n")
		handles = append(handles, fmt.Sprintf("%s = l.svcCtx.Models.%s.%s(l.ctx,data)", logicHandle.Result, modelStruct, logicHandle.Function))
	} else {
		handles = append(handles, fmt.Sprintf("%s = l.svcCtx.Models.%s.%s(%s)", logicHandle.Result, modelStruct, logicHandle.Function, logicHandle.Params))
	}
	handles = append(handles, fmt.Sprintf("if %s != nil{\n\t\treturn nil,err\n\t}", logicHandle.Result))

	return strings.Join(handles, "\n\t")
}

func newLogicHandlerByMethod(model string, apiGroups []spec.Group) map[string]spec.LogicHandle {

	logicModelInterface := make(map[string]spec.LogicHandle, 0)

	for _, g := range apiGroups {
		for _, r := range g.Routes {
			logicName := getLogicName(r)
			if r.Path == "/" {
				switch r.Method {
				case "GET":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:    model,
						Function: "FindOne",
						Params:   "l.ctx,req.id",
						Result:   "resp,err",
						IsModel:  false,
					}
				case "POST":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:    model,
						Function: "Create",
						Params:   "l.ctx",
						Result:   "_,err",
						IsModel:  true,
					}
				case "PUT":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:    model,
						Function: "Create",
						Params:   "l.ctx",
						Result:   "_,err",
						IsModel:  true,
					}
				case "DELETE":
					logicModelInterface[logicName] = spec.LogicHandle{
						Model:    model,
						Function: "Create",
						Params:   "l.ctx,req.id",
						Result:   "_,err",
						IsModel:  false,
					}
				}
			}
		}
	}

	return logicModelInterface
}
