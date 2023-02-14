package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"simple-main/http-rcp/pkg/utils/etcd"
)

/*
 @Author: 71made
 @Date: 2023/02/15 00:03
 @ProductName: init.go
 @Description:
*/

func InitClientConn(schema, serverName string) (*grpc.ClientConn, error) {

	rsv, err := etcd.BuildGRPCResolver()
	if err != nil {
		return nil, fmt.Errorf("can't create etcd client: %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", schema, serverName),
		grpc.WithResolvers(rsv), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		return nil, fmt.Errorf("can't connect server: %v", err)
	}
	return conn, nil
}
