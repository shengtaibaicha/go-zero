package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero/apps/api/community/internal/config"
	"go-zero/apps/rpc/file/client/fileupload"
	"go-zero/apps/rpc/user/client/login"
	"go-zero/apps/rpc/user/client/register"
	"go-zero/common/middleware"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config            config.Config
	CorsMiddleware    rest.Middleware
	JwtAuthMiddleware rest.Middleware
	FileClient        fileupload.FileUpload
	LoginClient       login.Login
	RegisterClient    register.Register
	RedisClient       *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		CorsMiddleware:    middleware.CorsMiddleware(),
		JwtAuthMiddleware: middleware.JwtAuthMiddleware(middleware.JwtAuthConfig{SecretKey: c.Jwt.SecretKey}),
		FileClient: fileupload.NewFileUpload(zrpc.MustNewClient(c.FileRpc,
			zrpc.WithDialOption(grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(30*1024*1024), // 客户端接收最大消息大小
				grpc.MaxCallSendMsgSize(30*1024*1024), // 客户端发送最大消息大小
			)),
		)),
		LoginClient:    login.NewLogin(zrpc.MustNewClient(c.UserRpc)),
		RegisterClient: register.NewRegister(zrpc.MustNewClient(c.UserRpc)),
		RedisClient:    redis.MustNewRedis(c.Redis),
	}
}
