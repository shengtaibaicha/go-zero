package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero/apps/user/rpc/internal/config"
	"go-zero/common/db"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	MDB    *gorm.DB
	Jwt    config.Jwt
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.RedisConf),
		MDB:    db.Init(c.Mysql.DataSource),
		Jwt:    c.Jwt,
	}
}
