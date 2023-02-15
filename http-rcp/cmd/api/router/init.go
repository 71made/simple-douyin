package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"simple-main/http-rcp/cmd/api/router/jwt"
)

/*
 @Author: 71made
 @Date: 2023/02/15 11:43
 @ProductName: init.go
 @Description:
*/

func Init(h *server.Hertz) {
	jwt.Init()
	register(h)
}
