package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/cmd/biz/core/service"
	"strconv"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/01/28 16:57
 @ProductName: feed.go
 @Description:
*/

var feedService = service.GetFeedServiceImpl()

func Feed(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	hlog.Info(token)

	lastTime, _ := strconv.ParseInt(c.Query("last_time"), 10, 64)

	resp := feedService.GetFeed(ctx, time.Unix(lastTime, 0), 1)
	c.JSON(http.StatusOK, resp)
}
