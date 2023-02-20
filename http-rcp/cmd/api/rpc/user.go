package rpc

import (
	"context"
	"simple-main/http-rcp/cmd/api/biz"
	rpc "simple-main/http-rcp/grpc_gen"
	usvr "simple-main/http-rcp/grpc_gen/user"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils/grpc"
)

/*
 @Author: 71made
 @Date: 2023/02/14 01:00
 @ProductName: user.go
 @Description:
*/

func userManagementClient() (client usvr.UserManagementClient, err error) {
	conn, err := grpc.InitClientConn(configs.Etcd, configs.UserServerName)
	if err != nil {
		return nil, err
	}
	return usvr.NewUserManagementClient(conn), nil
}

func CheckLoginUser(ctx context.Context, username, password string) (int64, *biz.Response) {

	req := &usvr.CheckLoginUserRequest{
		Username: username,
		Password: password,
	}

	client, err := userManagementClient()
	if err != nil {
		return biz.NotLoginUserId, biz.NewErrorResponse(err)
	}

	resp, err := client.CheckLoginUser(ctx, req)
	if err != nil {
		return biz.NotLoginUserId, biz.NewErrorResponse(err)
	}
	if resp != nil && resp.BaseResponse.StatusCode != rpc.Status_OK {
		return biz.NotLoginUserId, NewBizResponse(resp.BaseResponse)
	}

	return resp.UserId, biz.NewSuccessResponse(resp.BaseResponse.StatusMsg)
}

func CreateUser(ctx context.Context, username, password, avatar string) (*usvr.User, *biz.Response) {

	req := &usvr.CreateUserRequest{
		Username: username,
		Password: password,
		Avatar:   avatar,
	}

	client, err := userManagementClient()
	if err != nil {
		return nil, biz.NewErrorResponse(err)
	}

	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		return nil, biz.NewErrorResponse(err)
	}

	if resp != nil && resp.BaseResponse.StatusCode != rpc.Status_OK {
		return nil, NewBizResponse(resp.BaseResponse)
	}
	return resp.User, biz.NewSuccessResponse(resp.BaseResponse.StatusMsg)
}

func QueryUser(ctx context.Context, userId int64) (*usvr.User, *biz.Response) {
	req := &usvr.QueryUserRequest{UserId: userId}

	client, err := userManagementClient()
	if err != nil {
		return nil, biz.NewErrorResponse(err)
	}

	resp, err := client.QueryUser(ctx, req)
	if err != nil {
		return nil, biz.NewErrorResponse(err)
	}

	if resp != nil && resp.BaseResponse.StatusCode != rpc.Status_OK {
		return nil, NewBizResponse(resp.BaseResponse)
	}
	return resp.User, biz.NewSuccessResponse(resp.BaseResponse.StatusMsg)
}

func QueryUsers(ctx context.Context, userIds []int64) ([]*usvr.User, *biz.Response) {
	req := &usvr.QueryUsersRequest{UserIds: userIds}

	client, err := userManagementClient()
	if err != nil {
		return nil, biz.NewErrorResponse(err)
	}

	resp, err := client.QueryUsers(ctx, req)
	if err != nil {
		return nil, biz.NewErrorResponse(err)
	}

	if resp != nil && resp.BaseResponse.StatusCode != rpc.Status_OK {
		return nil, NewBizResponse(resp.BaseResponse)
	}
	return resp.UserList, biz.NewSuccessResponse(resp.BaseResponse.StatusMsg)
}
