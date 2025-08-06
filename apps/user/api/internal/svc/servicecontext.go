package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero/apps/user/api/internal/config"
	"go-zero/apps/user/rpc/client/login"
	"go-zero/apps/user/rpc/client/register"
)

type ServiceContext struct {
	Config   config.Config
	Redis    *redis.Redis
	Register register.Register
	Login    login.Login
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Redis:    redis.MustNewRedis(c.Redis),
		Register: register.NewRegister(zrpc.MustNewClient(c.UserRpc)),
		Login:    login.NewLogin(zrpc.MustNewClient(c.UserRpc)),
	}
}
