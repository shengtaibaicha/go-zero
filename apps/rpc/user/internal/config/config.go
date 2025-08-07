package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf redis.RedisConf
	Mysql     Mysql
	Jwt       Jwt
}

type Jwt struct {
	SecretKey string
}

type Mysql struct {
	DataSource string
}
