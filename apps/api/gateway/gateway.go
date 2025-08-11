package main

import (
	"flag"
	"fmt"
	"go-zero/apps/api/gateway/internal/config"
	"go-zero/apps/api/gateway/internal/handler"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/common/cors"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, cors.NewRestCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
