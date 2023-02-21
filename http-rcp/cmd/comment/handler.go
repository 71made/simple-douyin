package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"simple-main/http-rcp/cmd/comment/dal"
	"simple-main/http-rcp/cmd/comment/pack"
	csvr "simple-main/http-rcp/grpc_gen/comment"
)

/*
 @Author: 71made
 @Date: 2023/02/21 17:38
 @ProductName: handler.go
 @Description:
*/

func newServer(opts ...grpc.ServerOption) *grpc.Server {
	svr := grpc.NewServer(opts...)
	csvr.RegisterCommentManagementServer(svr, &CommentManagementServer{})
	return svr
}

type CommentManagementServer struct {
	csvr.UnimplementedCommentManagementServer
}

func (cms CommentManagementServer) CreateComment(ctx context.Context, req *csvr.CreateCommentRequest) (*csvr.CreateCommentResponse, error) {
	resp := &csvr.CreateCommentResponse{}

	newComment := &dal.Comment{
		UserId:  uint(req.UserId),
		VideoId: uint(req.VideoId),
		Content: req.Content,
	}

	if err := dal.CreateComment(ctx, newComment); err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.Comment = pack.Comment(newComment)
	resp.BaseResponse = pack.NewSuccessResponse("评论发表成功")
	return resp, nil
}

func (cms CommentManagementServer) DeleteComment(ctx context.Context, req *csvr.DeleteCommentRequest) (*csvr.DeleteCommentResponse, error) {
	resp := &csvr.DeleteCommentResponse{}

	comment, err := dal.QueryComment(ctx, req.Id)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	// 校验记录参数是否一致
	if comment == nil || comment.UserId != uint(req.UserId) || comment.VideoId != uint(req.VideoId) {
		resp.BaseResponse = pack.NewFailureResponse("参数或用户权限异常")
		return resp, nil
	}

	if err = dal.DeleteComment(ctx, req.Id, req.VideoId); err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.BaseResponse = pack.NewSuccessResponse("删除成功")
	return resp, nil
}

func (cms CommentManagementServer) QueryComments(ctx context.Context, req *csvr.QueryCommentsRequest) (*csvr.QueryCommentsResponse, error) {
	resp := &csvr.QueryCommentsResponse{}

	comments, err := dal.QueryComments(ctx, req.VideoId)
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		resp.CommentList = make([]*csvr.Comment, 0)
		return resp, nil
	}

	resp.CommentList = pack.Comments(comments)
	resp.BaseResponse = pack.NewSuccessResponse("获取成功")
	return resp, nil
}
