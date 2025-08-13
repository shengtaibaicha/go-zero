package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Minio        MinioConf
	Mysql        Mysql
	RedisExpires int
}

// MinioConf MinIO配置结构体，与配置文件中的Minio节点对应
type MinioConf struct {
	Endpoint  string // MinIO服务地址
	AccessKey string // 访问密钥
	SecretKey string // 密钥
	Bucket    string // 存储桶名称
	UseSSL    bool   // 是否使用SSL
	Expires   int64  // 预签名URL有效期（分钟）
}

type Mysql struct {
	DataSource string
}
