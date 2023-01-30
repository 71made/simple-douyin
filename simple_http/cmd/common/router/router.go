package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"net/http"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/core/controller"
)

/*
 @Author: 71made
 @Date: 2023/01/25 12:08
 @ProductName: router.go
 @Description:
*/

// Register 路由注册. 全部的路由注册配置都在此函数中
func Register(r *server.Hertz) {

	// 静态资源
	r.Static("/static", "./resources")

	root := r.Group("/douyin")
	// 获取视频流
	root.GET("/feed/", controller.Feed)
	{
		_user := root.Group("/user")
		// 用户信息
		_user.GET("/", controller.UserInfo)
		// 登陆
		_user.POST("/login/", controller.UserLogin)
		// 注册
		_user.POST("/register/", controller.UserRegister)

		_publish := root.Group("/publish")
		// 视频投稿
		_publish.POST("/action/", controller.Publish)
		// 获取视频列表
		_publish.GET("/list/", controller.PublishList)

		_favourite := root.Group("/favourite")
		// 视频点赞/取消点赞
		_favourite.POST("/action/", UnsupportedMethod)
		// 喜欢视频列表
		_favourite.GET("/list/", UnsupportedMethod)

		_comment := root.Group("/comment")
		// 发表评论
		_comment.POST("/action/", UnsupportedMethod)
		// 评论列表
		_comment.GET("/list/", UnsupportedMethod)

		_relation := root.Group("/relation")
		// 关注/取消关注
		_relation.POST("/action/", UnsupportedMethod)
		// 关注者列表
		_relation.GET("/follow/list/", UnsupportedMethod)
		// 粉丝列表
		_relation.GET("/follower/list/", UnsupportedMethod)
		// 好友列表
		_relation.GET("/friend/list/", UnsupportedMethod)

		_message := root.Group("/message")
		_message.GET("/chat/", UnsupportedMethod)
		_message.POST("/action/", UnsupportedMethod)
	}
}

func UnsupportedMethod(_ context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, biz.NewFailureResponse("暂不支持该接口服务"))
}
