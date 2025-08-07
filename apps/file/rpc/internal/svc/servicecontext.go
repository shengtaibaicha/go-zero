package svc

import (
	"github.com/minio/minio-go/v7"
	"go-zero/apps/file/rpc/internal/config"
	newminio "go-zero/apps/file/rpc/tools/minio"
	"go-zero/common/db"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	MinioClient *minio.Client // MinIO客户端实例
	MDB         *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		MinioClient: newminio.NewMinioClient(c),
		MDB:         db.Init(c.Mysql.DataSource),
	}
}
