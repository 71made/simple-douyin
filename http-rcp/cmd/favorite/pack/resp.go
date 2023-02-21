package pack

import (
	rpc "simple-main/http-rcp/grpc_gen"
	"time"
)

/*
@Author: 71made
@Date: 2023/02/15 12:18
@ProductName: resp.go
@Description:
*/

func NewErrorResponse(err error) *rpc.BaseResponse {
	return NewResponse(rpc.Status_Error, err.Error())
}

func NewSuccessResponse(successMsg string) *rpc.BaseResponse {
	return NewResponse(rpc.Status_OK, successMsg)
}

func NewFailureResponse(failureMsg string) *rpc.BaseResponse {
	return NewResponse(rpc.Status_Failure, failureMsg)
}

func NewResponse(code rpc.Status, msg string) *rpc.BaseResponse {
	return &rpc.BaseResponse{
		StatusCode:    code,
		StatusMsg:     msg,
		RespTimestamp: time.Now().Unix(),
	}
}
