package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"simple-main/cmd/common"
)

/*
 @Author: 71made
 @Date: 2023/01/24 22:13
 @ProductName: main.go
 @Description:
*/

func main() {
	tracer, cfg := tracing.NewServerTracer()

	h := server.New(
		server.WithHostPorts(":8080"),
		server.WithMaxRequestBodySize(50*1024*1024),
		//server.WithTransport(standard.NewTransporter),
		tracer,
	)
	common.Init(context.Background(), h)

	// use otel mw
	h.Use(tracing.ServerMiddleware(cfg))

	h.Spin()
}
