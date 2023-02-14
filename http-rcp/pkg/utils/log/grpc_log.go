package log

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

/*
 @Author: 71made
 @Date: 2023/02/13 18:48
 @ProductName: grpc_log.go
 @Description:
*/

type GrpcLogger struct {
	level  int
	logger *logging.Logger
}

// NewGrpcLogger 创建封装了 zap 的对象，该对象是对 LoggerV2 接口的实现
func NewGrpcLogger(logger *logging.Logger) *GrpcLogger {
	return &GrpcLogger{
		level:  99, // 日志等级, 默认最高等级 99, 所有日志信息都将打印
		logger: logger,
	}
}

// Info returns
func (gl *GrpcLogger) Info(args ...interface{}) {
	gl.logger.Info(args)
}

// Infoln returns
func (gl *GrpcLogger) Infoln(args ...interface{}) {
	gl.logger.Info(args...)
}

// Infof returns
func (gl *GrpcLogger) Infof(format string, args ...interface{}) {
	gl.logger.Infof(format, args...)
}

// Warning returns
func (gl *GrpcLogger) Warning(args ...interface{}) {
	gl.logger.Warn(args...)
}

// Warningln returns
func (gl *GrpcLogger) Warningln(args ...interface{}) {
	gl.logger.Warn(args...)
}

// Warningf returns
func (gl *GrpcLogger) Warningf(format string, args ...interface{}) {
	gl.logger.Warnf(format, args...)
}

// Error returns
func (gl *GrpcLogger) Error(args ...interface{}) {
	gl.logger.Error(args...)
}

// Errorln returns
func (gl *GrpcLogger) Errorln(args ...interface{}) {
	gl.logger.Error(args...)
}

// Errorf returns
func (gl *GrpcLogger) Errorf(format string, args ...interface{}) {
	gl.logger.Errorf(format, args...)
}

// Fatal returns
func (gl *GrpcLogger) Fatal(args ...interface{}) {
	gl.logger.Fatal(args...)
}

// Fatalln returns
func (gl *GrpcLogger) Fatalln(args ...interface{}) {
	gl.logger.Fatal(args...)
}

// Fatalf logs to fatal level
func (gl *GrpcLogger) Fatalf(format string, args ...interface{}) {
	gl.logger.Fatalf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (gl *GrpcLogger) V(v int) bool {
	return v <= gl.level
}
