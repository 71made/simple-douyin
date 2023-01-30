package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"path"
	"path/filepath"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/core/service"
	"simple-main/cmd/common/minio"
	"sync"
)

/*
 @Author: 71made
 @Date: 2023/01/27 19:57
 @ProductName: publish.go
 @Description: 发布视频、发布列表服务接口处理的 controller
*/

var publishService = service.GetPublishServiceImpl()

func Publish(ctx context.Context, c *app.RequestContext) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	hlog.Info("token :" + token)

	data, err := c.FormFile("data")
	if err != nil {
		resp := biz.NewErrorResponse(err)
		c.JSON(http.StatusOK, resp)
		return
	}

	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", 1, title+path.Ext(fileName))

	resp := publishService.PublishVideo(
		ctx,
		&service.VideoPublishRequest{
			VideoFinalName: finalName,
			UserId:         1,
			Title:          title,
		},
		func(dstPath string) (err error) {
			// wait group 控制并发
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				err = c.SaveUploadedFile(data, dstPath)
				wg.Done()
			}()

			var uploadErr error
			go func() {
				uploadErr = minio.UploadVideo(ctx, finalName, data)
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

func PublishList(ctx context.Context, c *app.RequestContext) {
	resp := publishService.GetPublishList(ctx, 1)
	c.JSON(http.StatusOK, resp)
}
