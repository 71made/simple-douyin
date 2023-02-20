package core

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/service/core"
	"simple-main/http-rcp/pkg/configs"
	"strconv"
)

/*
 @Author: 71made
 @Date: 2023/02/19 02:18
 @ProductName: feed.go
 @Description:
*/

var feedService = core.FeedServiceImpl()

// Feed
// @router /douyin/feed/ [POST]
func Feed(ctx context.Context, c *app.RequestContext) {
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	var userId int64
	if v != nil {
		userId = v.(*biz.User).Id
	} else {
		userId = biz.NotLoginUserId
	}

	lastTime, _ := strconv.ParseInt(c.Query("last_time"), 10, 64)

	resp := feedService.GetFeed(ctx, lastTime, userId)
	c.JSON(http.StatusOK, resp)
}
