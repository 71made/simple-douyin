package core

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/service/core"
	"simple-main/cmd/configs"
	"sync"
)

/*
 @Author: 71made
 @Date: 2023/01/27 19:57
 @ProductName: publish.go
 @Description: 发布视频、发布列表服务接口处理的 controller
*/

var publishService = core.GetPublishServiceImpl()

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
		resp := biz.NewErrorResponse(err)
		c.JSON(http.StatusOK, resp)
		return
	}

	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, title+path.Ext(fileName))

	resp := publishService.PublishVideo(
		ctx,
		&core.VideoPublishRequest{
			VideoFinalName: finalName,
			UserId:         1,
			Title:          title,
		},
		// 并发实现服务器保存和上传 MinIO
		func(dstPath string, uploadToMinIO func(data *multipart.FileHeader) error) (err error) {
			// wait group 控制并发
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				err = c.SaveUploadedFile(data, dstPath)
				wg.Done()
			}()

			var uploadErr error
			go func() {
				uploadErr = uploadToMinIO(data)
				wg.Done()
			}()

			wg.Wait()
			if err == nil {
				err = uploadErr
			}
			return
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
