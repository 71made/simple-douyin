package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/grpclog"
	"net"
	"simple-main/http-rcp/cmd/video/dal/db"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils"
	"simple-main/http-rcp/pkg/utils/etcd"
	"simple-main/http-rcp/pkg/utils/log"
	"simple-main/http-rcp/pkg/utils/minio"
)

/*
 @Author: 71made
 @Date: 2023/02/17 22:30
 @ProductName: main.go
 @Description:
*/

func preInit(ctx context.Context) {
	log.Init()
	db.Init()
	minio.Init(ctx)
}

func main() {
	preInit(context.Background())

	var port int
	flag.IntVar(&port, "port", configs.VideoServerPort, "port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", configs.ServerIP, port)

	listen, err := net.Listen(configs.TCP, addr)
	if err != nil {
		panic(err)
	}

	etcdClient, err := etcd.Register(configs.VideoServerName, addr)
	if err != nil {
		grpclog.Fatal("Video grpc server register to ETCD failed: ", err)
	}

	utils.DealSignal(func() {
		// 注销注册
		_ = etcd.Unregister(etcdClient, configs.VideoServerName, addr)
	})

	svr := newServer()
	grpclog.Info("Running video grpc server...")
	err = svr.Serve(listen)
	if err != nil {
		grpclog.Fatal("Video grpc server start failed: ", err)
		return
	}
}
