package common

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/simple-http/pkg/common/db"
	"simple-main/simple-http/pkg/common/jwt"
	"simple-main/simple-http/pkg/common/minio"
	"simple-main/simple-http/pkg/common/router"
)

/*
 @Author: 71made
 @Date: 2023/01/24 21:28
 @ProductName: init.go
 @Description: 集成初始化操作, 包括 db 数据库连接、hertz 路由注册等
*/

func Init(ctx context.Context, h *server.Hertz) {
	db.Init()
	minio.Init(ctx)
	jwt.Init()
	router.Register(h)

	// hlog init
	// 配置 hertz 日志格式, 默认格式
	hlog.SetLogger(hlog.DefaultLogger())
	hlog.SetLevel(hlog.LevelInfo)

}
