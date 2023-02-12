package second

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/biz/service/extra/second"
	"simple-main/simple-http/cmd/configs"
	"strconv"
)

/*
 @Author: 71made
 @Date: 2023/02/02 15:28
 @ProductName: relation.go
 @Description:
*/

var relationService = second.RelationServiceImpl()

// RelationAction
// @router /douyin/relation/action [POST]
func RelationAction(ctx context.Context, c *app.RequestContext) {
	req := &second.RelationActionRequest{}

	err := c.BindAndValidate(req)
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数绑定错误")))
		return
	}

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id
	req.UserId = userId

	resp := relationService.Action(ctx, req)
	c.JSON(http.StatusOK, resp)
}

// GetFollowerList
// @router /douyin/relation/follower/list [GET]
func GetFollowerList(ctx context.Context, c *app.RequestContext) {

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		resp := biz.NewFailureResponse("参数类型转换失败")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var thisUserId int64
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	if v != nil {
		thisUserId = v.(*biz.User).Id
	} else {
		thisUserId = biz.NotLoginUserId
	}

	resp := relationService.FollowerList(ctx, userId, thisUserId)
	c.JSON(http.StatusOK, resp)
}

// GetFollowList
// @router /douyin/relation/follow/list [GET]
func GetFollowList(ctx context.Context, c *app.RequestContext) {

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		resp := biz.NewFailureResponse("参数类型转换失败")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var thisUserId int64
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	if v != nil {
		thisUserId = v.(*biz.User).Id
	} else {
		thisUserId = biz.NotLoginUserId
	}

	resp := relationService.FollowList(ctx, userId, thisUserId)
	c.JSON(http.StatusOK, resp)
}

// GetFriendList
// @router /douyin/relation/friend/list [GET]
func GetFriendList(ctx context.Context, c *app.RequestContext) {

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		resp := biz.NewFailureResponse("参数类型转换失败")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := relationService.FriendList(ctx, userId)
	c.JSON(http.StatusOK, resp)
}
