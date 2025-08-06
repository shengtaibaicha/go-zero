package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero/apps/file/api/internal/config"
	"go-zero/apps/file/rpc/client/fileupload"
)

type ServiceContext struct {
	Config     config.Config
	FileClient fileupload.FileUpload
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		FileClient: fileupload.NewFileUpload(zrpc.MustNewClient(c.FileClient)),
	}
}
