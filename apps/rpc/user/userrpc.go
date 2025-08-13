package main

import (
	"flag"
	"fmt"

	"go-zero/apps/rpc/user/internal/config"
	loginServer "go-zero/apps/rpc/user/internal/server/login"
	otherServer "go-zero/apps/rpc/user/internal/server/other"
	registerServer "go-zero/apps/rpc/user/internal/server/register"
	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/userrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterLoginServer(grpcServer, loginServer.NewLoginServer(ctx))
		user.RegisterRegisterServer(grpcServer, registerServer.NewRegisterServer(ctx))
		user.RegisterOtherServer(grpcServer, otherServer.NewOtherServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
