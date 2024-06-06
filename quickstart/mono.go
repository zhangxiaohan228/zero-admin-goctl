package quickstart

import (
	_ "embed"
	"github.com/zhangxiaohan228/zero-admin-goctl/api/gogen"
	"github.com/zhangxiaohan228/zero-admin-goctl/pkg/golang"
	"github.com/zhangxiaohan228/zero-admin-goctl/util"
	"github.com/zhangxiaohan228/zero-admin-goctl/util/pathx"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	//go:embed idl/greet.api
	apiContent string
	//go:embed idl/svc.tpl
	svcContent string
	//go:embed idl/apilogic.tpl
	apiLogicContent string
	//go:embed idl/api.yaml
	apiEtcContent string

	apiWorkDir string
	rpcWorkDir string
)

func initAPIFlags() error {
	rpcWorkDir = filepath.Join(projectDir, "rpc")
	apiWorkDir = filepath.Join(projectDir, "api")
	if err := pathx.MkdirIfNotExist(apiWorkDir); err != nil {
		return err
	}

	apiFilename := filepath.Join(apiWorkDir, "greet.api")
	apiBytes := []byte(apiContent)
	if err := os.WriteFile(apiFilename, apiBytes, 0o666); err != nil {
		return err
	}

	gogen.VarStringDir = apiWorkDir
	gogen.VarStringAPI = apiFilename
	return nil
}

type mono struct {
	callRPC bool
}

func newMonoService(callRPC bool) mono {
	m := mono{callRPC}
	m.createAPIProject()
	return m
}

func (m mono) createAPIProject() {
	logx.Must(initAPIFlags())
	logx.Debug(">> Generating quickstart api project...")
	logx.Must(gogen.GoCommand(nil, nil))
	etcFile := filepath.Join(apiWorkDir, "etc", "greet.yaml")
	logx.Must(os.WriteFile(etcFile, []byte(apiEtcContent), 0o666))
	logicFile := filepath.Join(apiWorkDir, "internal", "logic", "pinglogic.go")
	svcFile := filepath.Join(apiWorkDir, "internal", "svc", "servicecontext.go")
	configPath := filepath.Join(apiWorkDir, "internal", "config")
	svcPath := filepath.Join(apiWorkDir, "internal", "svc")
	typesPath := filepath.Join(apiWorkDir, "internal", "types")
	svcPkg, err := golang.GetParentPackage(svcPath)
	logx.Must(err)
	typesPkg, err := golang.GetParentPackage(typesPath)
	logx.Must(err)
	configPkg, err := golang.GetParentPackage(configPath)
	logx.Must(err)

	var rpcClientPkg string
	if m.callRPC {
		rpcClientPath := filepath.Join(rpcWorkDir, "greet")
		rpcClientPkg, err = golang.GetParentPackage(rpcClientPath)
		logx.Must(err)
	}

	logx.Must(util.With("logic").Parse(apiLogicContent).SaveTo(map[string]any{
		"svcPkg":       svcPkg,
		"typesPkg":     typesPkg,
		"rpcClientPkg": rpcClientPkg,
		"callRPC":      m.callRPC,
	}, logicFile, true))

	logx.Must(util.With("svc").Parse(svcContent).SaveTo(map[string]any{
		"rpcClientPkg": rpcClientPkg,
		"configPkg":    configPkg,
		"callRPC":      m.callRPC,
	}, svcFile, true))
}

func (m mono) start() {
	if !m.callRPC {
		goModTidy(projectDir)
	}
	logx.Debug(">> Ready to start an API server...")
	logx.Debug(">> Run 'curl http://127.0.0.1:8888/ping' after service startup...")
	goStart(apiWorkDir)
}