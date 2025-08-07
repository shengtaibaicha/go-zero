package main

import (
	"flag"
	"fmt"
	"go-zero/common/middleware"

	"go-zero/apps/file/api/internal/config"
	"go-zero/apps/file/api/internal/handler"
	"go-zero/apps/file/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/upload-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	// jwt中间件
	authConfig := middleware.JwtAuthConfig{SecretKey: c.Jwt.SecretKey}
	server.Use(middleware.NewJwtAuthMiddleware(authConfig))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
