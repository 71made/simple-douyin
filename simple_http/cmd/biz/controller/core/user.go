package core

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/biz/service/core"
	"simple-main/simple-http/cmd/model"
	"simple-main/simple-http/pkg/configs"
	"strconv"
)

/*
 @Author: 71made
 @Date: 2023/01/25 13:50
 @ProductName: user.go
 @Description: 用户相关服务接口处理的 controller
*/

var userService = core.UserServiceImpl()

// UserInfo
// @router /douyin/user/ [GET]
func UserInfo(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp := biz.NewFailureResponse("请求参数异常")
		hlog.Errorf(fmt.Sprintf("msg : %s\n error: %v", "user_id 类型转换错误", err))
		c.JSON(http.StatusBadRequest, resp)
	}

	var thisUserId int64
	// 获取 JWT 回设的 userId
	v, _ := c.Get(configs.IdentityKey)
	if v != nil {
		thisUserId = v.(*biz.User).Id
	} else {
		thisUserId = biz.NotLoginUserId
	}

	resp := userService.UserInfo(ctx, userId, thisUserId)
	c.JSON(http.StatusOK, resp)
}

// UserLogin
// @router /douyin/user/login/ [POST]
// 注意: 使用 JWT 中间件后, 路由不再直接绑定此方法, 而是在 JWT 中配置的 Authenticator(认证器) 方法中调用,
// 而路由则直接绑定的 JWT 提供的 LoginHandler 方法. 对于 UserLogin 返回的 resp 和 status 调整为
// 通过 app.RequestContext 传递, 最终再通过 JWT 配置的 Unauthorized(未授权的)、LoginResponse 等回调方法写入 JSON 返回给客户端
func UserLogin(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	if len(username) == 0 || len(password) == 0 {
		c.Set("resp", biz.NewFailureResponse("用户名或密码不能为空"))
		c.Set("status", http.StatusBadRequest)
		return
	}

	resp := userService.Login(ctx, username, password)
	c.Set("resp", resp)
}

// UserRegister
// @router /douyin/user/register/ [POST]
func UserRegister(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, biz.NewFailureResponse("用户名或密码不能为空"))
		return
	}

	resp := userService.Register(
		ctx,
		&model.User{
			Username: username,
			Password: password,
		})

	c.JSON(http.StatusOK, resp)
}
