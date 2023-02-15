package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"simple-main/http-rcp/cmd/api/router"
	"simple-main/http-rcp/pkg/utils/minio"
)

/*
 @Author: 71made
 @Date: 2023/02/14 01:50
 @ProductName: main.go
 @Description:
*/

func perInit(ctx context.Context, h *server.Hertz) {
	router.Init(h)
	minio.Init(ctx)

	// hlog init
	// 配置 hertz 日志格式, 默认格式
	hlog.SetLogger(hlog.DefaultLogger())
	hlog.SetLevel(hlog.LevelInfo)
}

func main() {
	tracer, cfg := tracing.NewServerTracer()

	h := server.New(
		server.WithHostPorts(":8080"),
		server.WithMaxRequestBodySize(50*1024*1024),
		//server.WithTransport(standard.NewTransporter),
		tracer,
	)

	perInit(context.Background(), h)

	// use otel mw
	h.Use(tracing.ServerMiddleware(cfg))

	h.Spin()

}
