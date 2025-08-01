package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero/apps/user/api/internal/config"
	"go-zero/apps/user/rpc/client/register"
)

type ServiceContext struct {
	Config  config.Config
	Redis   *redis.Redis
	Regiter register.Register
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Redis:   redis.MustNewRedis(c.Redis),
		Regiter: register.NewRegister(zrpc.MustNewClient(c.UserRpc)),
	}
}
