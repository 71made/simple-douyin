package core

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"path"
	"path/filepath"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/service/core"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/18 22:36
 @ProductName: publish.go
 @Description:
*/

var publishService = core.PublishServiceImpl()

// Publish
// @router /douyin/publish/action/ [POST]
func Publish(ctx context.Context, c *app.RequestContext) {
	// 获取请求参数
	title := c.PostForm("title")
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	data, err := c.FormFile("data")
	if err != nil {
		hlog.Error(err)
		resp := biz.NewErrorResponse(err)
		c.JSON(http.StatusOK, resp)
		return
	}

	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, title+path.Ext(fileName))

	err = c.SaveUploadedFile(data, filepath.Join(configs.VideoPathPrefix, finalName))
	if err != nil {
		hlog.Error(err)
		resp := biz.NewErrorResponse(err)
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := publishService.PublishVideo(
		ctx,
		&core.VideoPublishRequest{
			VideoFinalName: finalName,
			UserId:         userId,
			Title:          title,
		})
	c.JSON(http.StatusOK, resp)
}

// PublishList
// @router /douyin/publish/list/ [GET]
func PublishList(ctx context.Context, c *app.RequestContext) {
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	resp := publishService.GetPublishList(ctx, userId)
	c.JSON(http.StatusOK, resp)
}
