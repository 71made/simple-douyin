package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"simple-main/http-rcp/cmd/favorite/dal"
	"simple-main/http-rcp/cmd/favorite/pack"
	fsvr "simple-main/http-rcp/grpc_gen/favorite"
)

/*
 @Author: 71made
 @Date: 2023/02/21 01:50
 @ProductName: handler.go
 @Description:
*/

func newServer(opts ...grpc.ServerOption) *grpc.Server {
	svr := grpc.NewServer(opts...)
	fsvr.RegisterFavoriteManagementServer(svr, &FavoriteManagementServer{})
	return svr
}

type FavoriteManagementServer struct {
	fsvr.UnimplementedFavoriteManagementServer
}

func (fms FavoriteManagementServer) Action(ctx context.Context, req *fsvr.FavoriteActionRequest) (*fsvr.FavoriteActionResponse, error) {
	resp := &fsvr.FavoriteActionResponse{}

	f := &dal.Favorite{
		UserId:  uint(req.UserId),
		VideoId: uint(req.VideoId),
	}

	switch req.Type {
	case fsvr.Action_Commit:
		{
			f.FavoriteType = dal.Favorable
		}
	case fsvr.Action_Cancel:
		{
			f.FavoriteType = dal.Unfavorable
		}
	case fsvr.Action_Unknown:
		fallthrough
	default:
		resp.BaseResponse = pack.NewFailureResponse("参数异常错误, 操作失败")
		return resp, nil
	}

	// 查找点赞记录
	found, err := dal.QueryFavorite(ctx, req.UserId, req.VideoId)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	// 没有记录则创建, 有则更新
	if found == nil {
		err = dal.CreateFavorite(ctx, f)
	} else if found.FavoriteType != f.GetFavoriteType() {
		// 并且只对于不同的 type, 才触发更新
		err = dal.UpdateFavorite(ctx, f)
	}

	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}
	resp.BaseResponse = pack.NewSuccessResponse("操作成功")
	return resp, nil
}

func (fms FavoriteManagementServer) QueryFavorites(ctx context.Context, req *fsvr.QueryFavoritesRequest) (*fsvr.QueryFavoritesResponse, error) {

	// 内部抽取处理函数
	packQueryRes := func(favorites []*dal.Favorite, err error) *fsvr.QueryFavoritesResponse {
		resp := &fsvr.QueryFavoritesResponse{}
		if err != nil {
			grpclog.Error(err)
			resp.BaseResponse = pack.NewErrorResponse(err)
			resp.FavoriteList = make([]*fsvr.Favorite, 0)
			return resp
		}

		resp.FavoriteList = pack.Favorites(favorites)
		resp.BaseResponse = pack.NewSuccessResponse("获取成功")
		return resp
	}

	if req.VideoIds == nil || len(req.VideoIds) == 0 {
		favorites, err := dal.QueryUserFavorites(ctx, req.UserId)
		return packQueryRes(favorites, err), nil
	} else {
		favorites, err := dal.QueryFavorites(ctx, req.UserId, req.VideoIds)
		return packQueryRes(favorites, err), nil
	}
}
