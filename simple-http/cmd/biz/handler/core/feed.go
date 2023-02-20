package core

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/biz/service/core"
	"simple-main/simple-http/pkg/configs"
	"strconv"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/01/28 16:57
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

	resp := feedService.GetFeed(ctx, time.Unix(lastTime, 0), userId)
	c.JSON(http.StatusOK, resp)
}
