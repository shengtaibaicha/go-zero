package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero/apps/common"
	"go-zero/apps/user/rpc/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	MDB    *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.RedisConf),
		MDB:    common.Init(c.Mysql.DataSource),
	}
}
