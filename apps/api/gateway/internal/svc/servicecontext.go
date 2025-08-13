package svc

import (
	"go-zero/apps/api/gateway/internal/config"
	"go-zero/apps/rpc/file/client/file"
	"go-zero/apps/rpc/tag/client/tag"
	"go-zero/apps/rpc/user/client/login"
	"go-zero/apps/rpc/user/client/other"
	"go-zero/apps/rpc/user/client/register"
	middleware2 "go-zero/common/middleware"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config            config.Config
	JwtAuthMiddleware rest.Middleware
	FileClient        file.File
	LoginClient       login.Login
	RegisterClient    register.Register
	RedisClient       *redis.Redis
	TagClient         tag.Tag
	OtherClient       other.Other
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		FileClient: file.NewFile(zrpc.MustNewClient(c.FileRpc,
			zrpc.WithDialOption(grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(30*1024*1024), // 客户端接收最大消息大小
				grpc.MaxCallSendMsgSize(30*1024*1024), // 客户端发送最大消息大小
			)),
		)),
		LoginClient:       login.NewLogin(zrpc.MustNewClient(c.UserRpc)),
		TagClient:         tag.NewTag(zrpc.MustNewClient(c.TagRpc)),
		RegisterClient:    register.NewRegister(zrpc.MustNewClient(c.UserRpc)),
		OtherClient:       other.NewOther(zrpc.MustNewClient(c.UserRpc)),
		RedisClient:       redis.MustNewRedis(c.Redis),
		JwtAuthMiddleware: middleware2.JwtAuthMiddleware(middleware2.JwtAuthConfig{SecretKey: c.Jwt.SecretKey, RedisClient: redis.MustNewRedis(c.Redis)}),
	}
}
