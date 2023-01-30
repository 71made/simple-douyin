package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/core/service"
	"simple-main/cmd/model"
	"strconv"
)

/*
 @Author: 71made
 @Date: 2023/01/25 13:50
 @ProductName: user.go
 @Description: 用户相关服务接口处理的 controller
*/

var userService = service.GetUserServiceImpl()

func UserInfo(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		resp := biz.NewFailureResponse("请求参数异常")
		hlog.Errorf(fmt.Sprintf("msg : %s\n error: %v", "user_id 类型转换错误", err))
		c.JSON(http.StatusBadRequest, resp)
	}
	token := c.Query("token")
	hlog.Info("token: ", token)

	resp := userService.UserInfo(ctx, userId)
	c.JSON(http.StatusOK, resp)
}

func UserLogin(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	if resp, ok := checkUserInfo(username, password); !ok {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := userService.Login(ctx, username, password)
	c.JSON(http.StatusOK, resp)
}

func UserRegister(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	if resp, ok := checkUserInfo(username, password); !ok {
		c.JSON(http.StatusBadRequest, resp)
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

// checkUserInfo 检查用户输入的用户名和密码字符串是否存在或为空
func checkUserInfo(username string, password string) (*biz.Response, bool) {
	if username == "" || password == "" {
		resp := biz.NewFailureResponse("用户名或密码不能为空")
		return resp, false
	}
	return nil, true
}
