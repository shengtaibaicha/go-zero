package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"go-zero/apps/user/api/internal/config"
	"go-zero/apps/user/api/internal/handler"
	"go-zero/apps/user/api/internal/svc"
	"go-zero/apps/user/api/middleware"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 解决跨域
	server := rest.MustNewServer(c.RestConf)
	server.Use(middleware.HeadersMiddleware())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
