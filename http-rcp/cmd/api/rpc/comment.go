package rpc

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/http-rcp/cmd/api/biz"
	rpc "simple-main/http-rcp/grpc_gen"
	csvr "simple-main/http-rcp/grpc_gen/comment"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils/grpc"
)

/*
 @Author: 71made
 @Date: 2023/02/21 19:20
 @ProductName: comment.go
 @Description:
*/

func commentManagementClient() (client csvr.CommentManagementClient, err error) {
	conn, err := grpc.InitClientConn(configs.Etcd, configs.CommentServerName)
	if err != nil {
		return nil, err
	}
	return csvr.NewCommentManagementClient(conn), nil
}

func CreateComment(ctx context.Context, videoId, userId int64, content string) (*csvr.Comment, *biz.Response) {
	req := &csvr.CreateCommentRequest{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}

	client, err := commentManagementClient()
	if err != nil {
		hlog.Error(err)
		return nil, biz.NewErrorResponse(err)
	}

	resp, err := client.CreateComment(ctx, req)
	if err != nil {
		hlog.Error(err)
		return nil, biz.NewErrorResponse(err)
	}

	return resp.Comment, NewBizResponse(resp.BaseResponse)
}

func DeleteComment(ctx context.Context, commentId, videoId, userId int64) *biz.Response {
	req := &csvr.DeleteCommentRequest{
		Id:      commentId,
		VideoId: videoId,
		UserId:  userId,
	}

	client, err := commentManagementClient()
	if err != nil {
		hlog.Error(err)
		return biz.NewErrorResponse(err)
	}

	resp, err := client.DeleteComment(ctx, req)
	if err != nil {
		hlog.Error(err)
		return biz.NewErrorResponse(err)
	}
	return NewBizResponse(resp.BaseResponse)
}

func QueryComments(ctx context.Context, videoId int64) ([]*csvr.Comment, *biz.Response) {
	req := &csvr.QueryCommentsRequest{VideoId: videoId}

	client, err := commentManagementClient()
	if err != nil {
		hlog.Error(err)
		return make([]*csvr.Comment, 0), biz.NewErrorResponse(err)
	}

	resp, err := client.QueryComments(ctx, req)
	if err != nil {
		hlog.Error(err)
		return make([]*csvr.Comment, 0), biz.NewErrorResponse(err)
	}
	if resp.BaseResponse != nil && resp.BaseResponse.StatusCode != rpc.Status_OK {
		return make([]*csvr.Comment, 0), NewBizResponse(resp.BaseResponse)
	}

	return resp.CommentList, NewBizResponse(resp.BaseResponse)
}
