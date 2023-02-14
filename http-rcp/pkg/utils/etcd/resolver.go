package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	gresolver "google.golang.org/grpc/resolver"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/14 23:52
 @ProductName: resolver.go
 @Description:
*/

func NewClient() (client *clientv3.Client, err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{configs.EtcdURL},
		DialTimeout: configs.EtcdDialTimeout,
	})
	return
}

func BuildGRPCResolver() (gresolver.Builder, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return resolver.NewBuilder(client)
}
