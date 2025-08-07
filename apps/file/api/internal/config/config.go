package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	FileRpc zrpc.RpcClientConf
	Jwt     Jwt
}

type Jwt struct {
	SecretKey string
}
