package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/01/29 21:05
 @ProductName: init.go
 @Description: 对接创建 MinIO 服务 Client
*/

var defaultClient *minio.Client

func Init(ctx context.Context) {
	// 构建 client
	var err error
	defaultClient, err = minio.New(configs.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(configs.AccessKeyID, configs.SecretAccessKey, ""),
		Secure: configs.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	// 检查目标桶是否存在
	found, err := defaultClient.BucketExists(ctx, configs.BucketName)
	if err != nil {
		panic(err)
	}
	// 不存在尝试创建
	if !found {
		if err = defaultClient.MakeBucket(ctx, configs.BucketName,
			minio.MakeBucketOptions{Region: "cn-north-1", ObjectLocking: true}); err != nil {
			panic("Make bucket fail, err: " + err.Error())
		}
	}
}

func Client() *minio.Client {
	if defaultClient == nil {
		panic("MinIO client doesn't init")
	}
	return defaultClient
}
