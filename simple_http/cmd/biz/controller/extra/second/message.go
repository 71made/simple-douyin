package second

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/service/extra/second"
	"simple-main/cmd/configs"
	"strconv"
)

/*
 @Author: 71made
 @Date: 2023/02/02 15:28
 @ProductName: message.go
 @Description:
*/

var messageService = second.MessageServiceImpl()

// MessageChat
// @router /douyin/message/chat [POST]
func MessageChat(ctx context.Context, c *app.RequestContext) {

	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		resp := biz.NewFailureResponse("参数类型转换失败")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	resp := messageService.Chat(ctx, userId, toUserId)
	c.JSON(http.StatusOK, resp)
}

// MessageAction
// @router /douyin/message/action [GET]
func MessageAction(ctx context.Context, c *app.RequestContext) {

	req := &second.MessageRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.Error(err)
		resp := biz.NewFailureResponse("参数绑定失败")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	req.UserId = userId
	resp := messageService.Action(ctx, req)
	c.JSON(http.StatusOK, resp)
}
