package main

import (
	"flag"
	"fmt"
	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/config"
	fileuploadServer "go-zero/apps/rpc/file/internal/server/fileupload"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/filerpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		file.RegisterFileUploadServer(grpcServer, fileuploadServer.NewFileUploadServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 设置grpc接收的最大文件为30M
	s.AddOptions(grpc.MaxRecvMsgSize(30 * 1024 * 1024))
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
