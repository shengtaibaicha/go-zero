package newminio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go-zero/apps/file/rpc/internal/config"
)

func NewMinioClient(c config.Config) *minio.Client {
	// 1. 初始化MinIO客户端
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
		Secure: c.Minio.UseSSL, // 是否使用HTTPS
	})
	if err != nil {
		panic(fmt.Sprintf("初始化MinIO客户端失败: %v", err))
	}

	// 2. 检查存储桶是否存在，不存在则创建
	exists, err := minioClient.BucketExists(context.Background(), c.Minio.Bucket)
	if err != nil {
		panic(fmt.Sprintf("检查存储桶是否存在失败: %v", err))
	}
	if !exists {
		// 创建存储桶，公开访问（根据需求调整权限）
		err = minioClient.MakeBucket(context.Background(), c.Minio.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			panic(fmt.Sprintf("创建存储桶失败: %v", err))
		}
	}
	return minioClient
}
