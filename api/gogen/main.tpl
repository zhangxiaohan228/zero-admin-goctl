package main

import (
	"flag"
	"fmt"
    "zero-admin/server/system/plugins/rouetrs"
	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterRouters(server, ctx, rouetrs.NewPluginRouter(server))

    if c.Log.Level == "debug" {
        httpx.SetErrorHandler(xerr.DebugErrorHandler)
    } else {
        httpx.SetErrorHandler(xerr.ErrorHandler)
    }

    httpx.SetOkHandler(xerr.OkHandler)
    logx.DisableStat() // 关闭系统硬件监控,线上建议打开

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
