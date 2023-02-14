package rpc

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"simple-main/http-rcp/cmd/api/biz"
	rpc "simple-main/http-rcp/grpc_gen"
	"simple-main/http-rcp/grpc_gen/user"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils/grpc"
)

/*
 @Author: 71made
 @Date: 2023/02/14 01:00
 @ProductName: user.go
 @Description:
*/

func userServiceClient() (client user.UserServiceClient, err error) {
	conn, err := grpc.InitClientConn(configs.Etcd, configs.UserServerName)
	if err != nil {
		return nil, err
	}
	return user.NewUserServiceClient(conn), nil
}

func CheckLoginUser(ctx context.Context, username, password string) (int64, error) {
	h := md5.New()
	if _, err := io.WriteString(h, password); err != nil {
		return biz.NotLoginUserId, err
	}

	// MD5 摘要算法处理密码
	password = fmt.Sprintf("%x", h.Sum(nil))

	req := &user.CheckLoginUserRequest{
		Username: username,
		Password: password,
	}

	client, err := userServiceClient()

	if err != nil {
		return biz.NotLoginUserId, err
	}

	resp, err := client.CheckLoginUser(ctx, req)
	if err != nil {
		return biz.NotLoginUserId, err
	}
	if resp != nil && resp.BaseResponse.GetStatusCode() != rpc.Status_OK {
		return biz.NotLoginUserId, errors.New("user rpc server has error")
	}

	fmt.Println(resp)

	return resp.UserId, nil
}
