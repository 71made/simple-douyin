package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc/grpclog"
	"net"
	"simple-main/http-rcp/cmd/user/dal/db"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils"
	"simple-main/http-rcp/pkg/utils/etcd"
	"simple-main/http-rcp/pkg/utils/log"
)

/*
 @Author: 71made
 @Date: 2023/02/13 02:00
 @ProductName: main.go
 @Description:
*/

func preInit() {
	log.Init()
	db.Init()
}

func main() {
	preInit()

	var port int
	flag.IntVar(&port, "port", configs.UserServerPort, "port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", configs.ServerIP, port)

	listen, err := net.Listen(configs.TCP, addr)
	if err != nil {
		panic(err)
	}

	etcdClient, err := etcd.Register(configs.UserServerName, addr)
	if err != nil {
		grpclog.Fatal("User grpc server register to ETCD failed: ", err)
	}

	utils.DealSignal(func() {
		// 注销注册
		_ = etcd.Unregister(etcdClient, configs.UserServerName, addr)
	})

	svr := newServer()
	grpclog.Info("Running user grpc server...")
	err = svr.Serve(listen)
	if err != nil {
		grpclog.Fatal("User grpc server start failed: ", err)
		return
	}
}
