package rpc

import (
	"context"
	"google.golang.org/grpc/grpclog"
	"simple-main/http-rcp/cmd/video/pack"
	rpc "simple-main/http-rcp/grpc_gen"
	fsvr "simple-main/http-rcp/grpc_gen/favorite"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils/grpc"
)

/*
 @Author: 71made
 @Date: 2023/02/21 03:06
 @ProductName: favorite.go
 @Description:
*/

func favoriteManagementClient() (client fsvr.FavoriteManagementClient, err error) {
	conn, err := grpc.InitClientConn(configs.Etcd, configs.FavoriteServerName)
	if err != nil {
		return nil, err
	}
	return fsvr.NewFavoriteManagementClient(conn), nil
}

func QueryUserFavorites(ctx context.Context, userId int64) ([]*fsvr.Favorite, *rpc.BaseResponse) {
	req := &fsvr.QueryFavoritesRequest{
		UserId: userId,
	}

	client, err := favoriteManagementClient()
	if err != nil {
		grpclog.Error(err)
		return make([]*fsvr.Favorite, 0), pack.NewErrorResponse(err)
	}

	resp, err := client.QueryFavorites(ctx, req)
	if err != nil {
		grpclog.Error(err)
		return make([]*fsvr.Favorite, 0), pack.NewErrorResponse(err)
	}

	if resp.BaseResponse != nil && resp.BaseResponse.StatusCode != rpc.Status_OK {
		return make([]*fsvr.Favorite, 0), resp.BaseResponse
	}

	return resp.FavoriteList, resp.BaseResponse

}
