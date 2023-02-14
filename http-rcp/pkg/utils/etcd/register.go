package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/14 22:34
 @ProductName: register.go
 @Description:
*/

func Register(serverName, addr string) (*clientv3.Client, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	em, err := endpoints.NewManager(client, serverName)
	if err != nil {
		return nil, err
	}
	// 创建租约
	lease, _ := client.Grant(context.TODO(), configs.EtcdTTL)
	// 挂载服务节点
	err = em.AddEndpoint(context.TODO(), fmt.Sprintf("%s/%s", serverName, addr),
		endpoints.Endpoint{Addr: addr}, clientv3.WithLease(lease.ID))
	if err != nil {
		return nil, err
	}
	// 开启保活机制
	alive, err := client.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			_ = <-alive
		}

	}()
	return client, nil
}
