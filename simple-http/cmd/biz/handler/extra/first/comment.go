package first

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/biz/service/extra/first"
	"simple-main/simple-http/pkg/configs"
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

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	// 构造 req
	req := &first.CommentRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数绑定失败")))
		return
	}
	req.UserId = userId

	if (req.ActionType == first.PublishComment && len(req.Content) == 0) ||
		(req.ActionType == first.RemoveComment && req.CommentId == 0) {
		hlog.Error("请求参数缺失")
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("请求参数缺失")))
		return
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

	var userId int64
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	if v != nil {
		userId = v.(*biz.User).Id
	} else {
		userId = biz.NotLoginUserId
	}

	resp := commentServiceImpl.CommentList(ctx, userId, videoId)
	c.JSON(http.StatusOK, resp)
}
