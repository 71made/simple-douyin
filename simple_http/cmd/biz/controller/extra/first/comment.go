package first

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/service/extra/first"
	"simple-main/cmd/configs"
	"strconv"
)

/*
 @Author: 71made
 @Date: 2023/01/31 15:42
 @ProductName: comment.go
 @Description:
*/

var commentServiceImpl = first.CommentServiceImpl()

// CommentAction
// @router /douyin/comment/action/ [POST]
func CommentAction(ctx context.Context, c *app.RequestContext) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数类型转换错误")))
		return
	}
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数类型转换错误")))
		return
	}
	commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	content := c.Query("comment_text")
	if (actionType == first.PublishComment && len(content) == 0) ||
		(actionType == first.RemoveComment && commentId == 0) {
		hlog.Error("请求参数缺失")
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("请求参数缺失")))
		return
	}

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	// 构造 req
	req := &first.CommentRequest{
		UserId:     userId,
		VideoId:    videoId,
		ActionType: actionType,
		Content:    content,
		CommentId:  commentId,
	}

	resp := commentServiceImpl.Action(ctx, req)
	c.JSON(http.StatusOK, resp)
}

// GetCommentList
// @router /douyin/comment/list/ [GET]
func GetCommentList(ctx context.Context, c *app.RequestContext) {
	var videoId int64
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数类型转换错误")))
		return
	}

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	resp := commentServiceImpl.CommentList(ctx, userId, videoId)
	c.JSON(http.StatusOK, resp)
}
