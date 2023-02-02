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
 @Date: 2023/01/31 15:37
 @ProductName: favorite.go
 @Description:
*/

var favoriteServiceImpl = first.FavoriteServiceImpl()

// FavoriteAction
// @router /douyin/favorite/action/ [POST]
func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	var req first.FavoriteRequest

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数类型转换错误")))
		return
	}

	acType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		hlog.Error(err)
		c.JSON(http.StatusBadRequest, biz.NewErrorResponse(fmt.Errorf("参数类型转换错误")))
		return
	}

	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	// 构建 req
	req.VideoId = videoId
	req.ActionType = uint(acType)
	req.UserId = userId

	resp := favoriteServiceImpl.Action(ctx, &req)
	c.JSON(http.StatusOK, resp)

}

// GetFavoriteList
// @route /douyin/favorite/list/ [GET]
func GetFavoriteList(ctx context.Context, c *app.RequestContext) {
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	userId := v.(*biz.User).Id

	resp := favoriteServiceImpl.FavoriteList(ctx, userId)
	c.JSON(http.StatusOK, resp)
}
