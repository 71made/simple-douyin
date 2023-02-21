package configs

import "time"

/*
 @Author: 71made
 @Date: 2023/02/14 22:43
 @ProductName: etcd.go
 @Description:
*/

const Etcd = "etcd"
const HTTP = "http"
const EtcdPort = ":2379"
const EtcdURL = ServerAddrPrefix + EtcdPort
const EtcdDialTimeout = 5 * time.Second
const EtcdTTL = 10

const (
	ServerNamePrefix = "douyin"
	ServerAddrPrefix = HTTP + "://" + ServerIP
)

const (
	UserServerName     = ServerNamePrefix + "/" + UserServer
	VideoServerName    = ServerNamePrefix + "/" + VideoServer
	FavoriteServerName = ServerNamePrefix + "/" + FavoriteServer
	CommentServerName  = ServerNamePrefix + "/" + CommentServer
)
