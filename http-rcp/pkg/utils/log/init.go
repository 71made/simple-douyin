package log

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"google.golang.org/grpc/grpclog"
)

/*
 @Author: 71made
 @Date: 2023/02/13 19:01
 @ProductName: init.go
 @Description:
*/

func Init() {
	var logger = logging.NewLogger("grpc-logger")
	grpclog.SetLoggerV2(NewGrpcLogger(logger))
}
