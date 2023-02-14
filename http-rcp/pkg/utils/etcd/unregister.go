package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

/*
 @Author: 71made
 @Date: 2023/02/14 23:20
 @ProductName: unregister.go
 @Description:
*/

func Unregister(client *clientv3.Client, serverName, addr string) error {

	if client != nil {
		em, err := endpoints.NewManager(client, serverName)
		if err != nil {
			return err
		}
		err = em.DeleteEndpoint(context.TODO(), fmt.Sprintf("%s/%s", serverName, addr))
		if err != nil {
			return err
		}
		return err
	}

	return nil
}
