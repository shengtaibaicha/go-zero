package svc

import (
	"go-zero/apps/rpc/file/internal/config"
	"go-zero/apps/rpc/file/tools/minio"
	"go-zero/common/db"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	MinioClient *minio.Client // MinIO客户端实例
	MDB         *gorm.DB
	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		MinioClient: newminio.NewMinioClient(c),
		MDB:         db.Init(c.Mysql.DataSource),
		RedisClient: redis.MustNewRedis(c.Redis.RedisConf),
	}
}
