package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io"
	"simple-main/http-rcp/cmd/user/dal"
	"simple-main/http-rcp/cmd/user/pack"
	usvr "simple-main/http-rcp/grpc_gen/user"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/13 02:00
 @ProductName: handler.go
 @Description:
*/

func newServer(opts ...grpc.ServerOption) *grpc.Server {
	svr := grpc.NewServer(opts...)
	usvr.RegisterUserManagementServer(svr, &UserManagementServerImpl{})
	return svr
}

// UserManagementServerImpl
// 继承 grpc 提供的类, 重写接口实现服务功能
type UserManagementServerImpl struct {
	usvr.UnimplementedUserManagementServer
}

func (ums UserManagementServerImpl) CheckLoginUser(ctx context.Context, req *usvr.CheckLoginUserRequest) (*usvr.CheckLoginUserResponse, error) {
	resp := &usvr.CheckLoginUserResponse{}

	username := req.Username
	password := req.Password
	h := md5.New()
	if _, err := io.WriteString(h, password); err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}
	// MD5 摘要算法处理密码
	password = fmt.Sprintf("%x", h.Sum(nil))

	user, err := dal.QueryUser(ctx, username)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	if user.Password != password {
		resp.BaseResponse = pack.NewFailureResponse("用户名或密码错误")
		return resp, nil
	}

	return &usvr.CheckLoginUserResponse{
		UserId:       int64(user.ID),
		BaseResponse: pack.NewSuccessResponse("校验通过"),
	}, nil
}

func (ums UserManagementServerImpl) CreateUser(ctx context.Context, req *usvr.CreateUserRequest) (*usvr.CreateUserResponse, error) {
	resp := &usvr.CreateUserResponse{}

	found, err := dal.IsExistUser(ctx, req.Username)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	if found {
		resp.BaseResponse = pack.NewFailureResponse("该用户名已被使用")
		return resp, nil
	}

	h := md5.New()
	if _, err = io.WriteString(h, req.Password); err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}
	// MD5 摘要算法处理密码
	password := fmt.Sprintf("%x", h.Sum(nil))

	// 构建 user
	user := &dal.User{
		Username: req.Username,
		Password: password,
		Avatar:   req.Avatar,
	}

	// 设置默认头像
	if len(user.Avatar) == 0 {
		user.Avatar = configs.EmptyAvatarName
	}

	if err = dal.CreateUser(ctx, user); err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.User = pack.User(user)
	resp.BaseResponse = pack.NewSuccessResponse("创建成功")
	return resp, nil
}

func (ums UserManagementServerImpl) QueryUsers(ctx context.Context, req *usvr.QueryUsersRequest) (*usvr.QueryUsersResponse, error) {
	resp := &usvr.QueryUsersResponse{}

	userIds := req.UserIds

	if len(userIds) == 0 {
		resp.UserList = pack.Users(nil)
		resp.BaseResponse = pack.NewSuccessResponse("获取成功")
		return resp, nil
	}

	users, err := dal.QueryUsersByIds(ctx, userIds)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.UserList = pack.Users(users)
	resp.BaseResponse = pack.NewSuccessResponse("获取成功")
	return resp, nil
}

func (ums UserManagementServerImpl) QueryUser(ctx context.Context, req *usvr.QueryUserRequest) (*usvr.QueryUserResponse, error) {
	resp := &usvr.QueryUserResponse{}

	userId := req.UserId

	user, err := dal.QueryUserById(ctx, userId)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.User = pack.User(user)
	resp.BaseResponse = pack.NewSuccessResponse("获取成功")
	return resp, nil
}
