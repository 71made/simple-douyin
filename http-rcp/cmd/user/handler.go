package main

import (
	"context"
	"google.golang.org/grpc"
	rpc "simple-main/http-rcp/grpc_gen"
	"simple-main/http-rcp/grpc_gen/user"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/02/13 02:00
 @ProductName: handler.go
 @Description:
*/

type UserServiceServerImpl struct {
	user.UnimplementedUserServiceServer
}

func (uss UserServiceServerImpl) CheckLoginUser(ctx context.Context, req *user.CheckLoginUserRequest) (*user.CheckLoginUserResponse, error) {
	return &user.CheckLoginUserResponse{
		UserId: 1,
		BaseResponse: &rpc.BaseResponse{
			StatusCode:    rpc.Status_OK,
			StatusMsg:     "请求成功",
			RespTimestamp: time.Now().Unix(),
		},
	}, nil
}

//func (uss UserServiceServerImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
//}
//func (uss UserServiceServerImpl) QueryUsers(ctx context.Context, req *user.QueryUsersRequest) (*user.QueryUsersResponse, error) {
//}
//func (uss UserServiceServerImpl) QueryUser(ctx context.Context, req *user.QueryUserRequest) (*user.QueryUserResponse, error) {
//}

func newServer(opts ...grpc.ServerOption) *grpc.Server {
	svr := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(svr, &UserServiceServerImpl{})
	return svr
}
