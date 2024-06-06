package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/zhangxiaohan228/zero-admin-goctl/api/spec"
	"github.com/zhangxiaohan228/zero-admin-goctl/config"
	"github.com/zhangxiaohan228/zero-admin-goctl/util/format"
	"github.com/zhangxiaohan228/zero-admin-goctl/util/pathx"
	"github.com/zhangxiaohan228/zero-admin-goctl/vars"
)

const contextFilename = "service_context"

//go:embed svc.tpl
var contextTemplate string

func genServiceContext(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, contextFilename)
	if err != nil {
		return err
	}

	var middlewareStr string
	var middlewareAssignment string
	middlewares := getMiddleware(api)
	model := getModel(api)
	ifMysql, modelSvc := genSqlConn(model)

	for _, item := range middlewares {
		middlewareStr += fmt.Sprintf("%s rest.Middleware\n", item)
		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		middlewareAssignment += fmt.Sprintf("%s: %s,\n", item,
			fmt.Sprintf("middleware.New%s().%s", strings.Title(name), "Handle"))
	}

	configImport := "\"" + pathx.JoinPackages(rootPkg, configDir) + "\""
	if len(middlewareStr) > 0 {
		configImport += "\n\t\"" + pathx.JoinPackages(rootPkg, middlewareDir) + "\""
		configImport += fmt.Sprintf("\n\t\"%s/rest\"", vars.ProjectOpenSourceURL)
	}
	if model != "" {
		configImport += "\n\t\"" + "zero-admin/server/model/model" + "\""
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          contextDir,
		filename:        filename + ".go",
		templateName:    "contextTemplate",
		category:        category,
		templateFile:    contextTemplateFile,
		builtinTemplate: contextTemplate,
		data: map[string]string{
			"configImport":         configImport,
			"config":               "config.Config",
			"middleware":           middlewareStr,
			"middlewareAssignment": middlewareAssignment,
			"ifMysql":              ifMysql,
			"models":               modelSvc,
		},
	})
}

func genSqlConn(apiModel string) (string, string) {
	if apiModel != "" {
		return "\tdbSourceStr := dbSource(c.DB)\n\tsqlConn := sqlx.NewMysql(dbSourceStr)", "\t\tModels:         model.NewModel(sqlConn),"
	}
	return "", ""
}
