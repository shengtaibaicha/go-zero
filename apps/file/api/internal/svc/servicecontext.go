package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero/apps/file/api/internal/config"
	"go-zero/apps/file/rpc/client/fileupload"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config     config.Config
	FileClient fileupload.FileUpload
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		FileClient: fileupload.NewFileUpload(zrpc.MustNewClient(c.FileRpc,
			zrpc.WithDialOption(grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(30*1024*1024), // 客户端接收最大消息大小
				grpc.MaxCallSendMsgSize(30*1024*1024), // 客户端发送最大消息大小
			)),
		)),
	}
}
