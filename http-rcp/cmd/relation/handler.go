package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"simple-main/http-rcp/cmd/relation/dal"
	"simple-main/http-rcp/cmd/relation/pack"
	rsvr "simple-main/http-rcp/grpc_gen/relation"
)

/*
 @Author: 71made
 @Date: 2023/02/22 22:58
 @ProductName: handler.go
 @Description:
*/

func newServer(opts ...grpc.ServerOption) *grpc.Server {
	svr := grpc.NewServer(opts...)
	rsvr.RegisterRelationManagementServer(svr, &RelationManagementServer{})
	return svr
}

type RelationManagementServer struct {
	rsvr.UnimplementedRelationManagementServer
}

func (rms RelationManagementServer) CreateRelation(ctx context.Context, req *rsvr.CreateRelationRequest) (*rsvr.CreateRelationResponse, error) {
	resp := &rsvr.CreateRelationResponse{}

	// 构建 relation, 默认创建的 relation 的为关注关系
	newRelation := &dal.Relation{
		UserId:       uint(req.UserId),
		FollowerId:   uint(req.FollowerId),
		FollowType:   uint(rsvr.Action_Following),
		FriendStatus: dal.Unknown,
	}

	if err := dal.CreateRelation(ctx, newRelation); err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.BaseResponse = pack.NewSuccessResponse("操作成功")
	return resp, nil
}

func (rms RelationManagementServer) UpdateRelation(ctx context.Context, req *rsvr.UpdateRelationRequest) (*rsvr.UpdateRelationResponse, error) {
	resp := &rsvr.UpdateRelationResponse{}

	relation, err := dal.QueryRelation(ctx, req.FollowerId, req.UserId)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	if uint(req.ActionType) == relation.FollowType {
		resp.BaseResponse = pack.NewSuccessResponse("重复操作")
		return resp, nil
	}

	switch req.ActionType {
	case rsvr.Action_Following:
		fallthrough
	case rsvr.Action_NotFollowing:
		{
			relation.FollowType = uint(req.ActionType)
			if err := dal.UpdateRelation(ctx, relation); err != nil {
				grpclog.Error(err)
				resp.BaseResponse = pack.NewErrorResponse(err)
				return resp, nil
			}

			resp.BaseResponse = pack.NewSuccessResponse("操作成功")
			return resp, nil
		}
	case rsvr.Action_Unknown:
		fallthrough
	default:
		{
			grpclog.Error(fmt.Sprint("update request's action type is not available."))
			resp.BaseResponse = pack.NewFailureResponse("参数异常错误")
			return resp, nil
		}
	}
}

func (rms RelationManagementServer) QueryRelation(ctx context.Context, req *rsvr.QueryRelationRequest) (*rsvr.QueryRelationResponse, error) {
	resp := &rsvr.QueryRelationResponse{}

	relation, err := dal.QueryRelation(ctx, req.ThisUserId, req.AnotherUserId)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.BaseResponse = pack.NewSuccessResponse("获取成功")
	resp.Relation = pack.Relation(relation)
	return resp, nil
}

func (rms RelationManagementServer) QueryRelations(ctx context.Context, req *rsvr.QueryRelationsRequest) (*rsvr.QueryRelationsResponse, error) {
	resp := &rsvr.QueryRelationsResponse{}

	var relations []*dal.Relation

	if req.ToUserIds != nil && len(req.ToUserIds) != 0 {
		// 查询该用户对所有列表用户的关注/粉丝关系
		var err error
		relations, err = dal.QueryRelations(ctx, req.ThisUserId, req.ToUserIds)
		if err != nil {
			grpclog.Error(err)
			resp.BaseResponse = pack.NewErrorResponse(err)
			return resp, nil
		}
	} else if req.FriendRelation {
		// 查询朋友关系
		var err error
		relations, err = dal.QueryFriendRelations(ctx, req.ThisUserId)
		if err != nil {
			grpclog.Error(err)
			resp.BaseResponse = pack.NewErrorResponse(err)
			return resp, nil
		}

	} else if req.FollowerRelation {
		// 查询粉丝关系
		var err error
		relations, err = dal.QueryFollowerRelations(ctx, req.ThisUserId)
		if err != nil {
			grpclog.Error(err)
			resp.BaseResponse = pack.NewErrorResponse(err)
			return resp, nil
		}
	} else if req.FollowRelation {
		// 查询关注关系
		var err error
		relations, err = dal.QueryFollowRelations(ctx, req.ThisUserId)
		if err != nil {
			grpclog.Error(err)
			resp.BaseResponse = pack.NewErrorResponse(err)
			return resp, nil
		}
	} else {
		grpclog.Error(fmt.Errorf("error query relations request, no result"))
		resp.BaseResponse = pack.NewFailureResponse("参数异常错误")
		resp.RelationList = pack.Relations(nil)
		return resp, nil
	}

	resp.BaseResponse = pack.NewSuccessResponse("获取成功")
	resp.RelationList = pack.Relations(relations)
	return resp, nil

}
