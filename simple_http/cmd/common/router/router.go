package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"net/http"
	"simple-main/cmd/biz"
	"simple-main/cmd/biz/core/controller"
	"simple-main/cmd/common/jwt"
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
	root.GET("/feed/", append([]app.HandlerFunc{func(ctx context.Context, c *app.RequestContext) {
		// 对于 Feed 接口, 如果传入了 token, 此处需要手动调用 JWT 的 mw 校验和解析 token
		if token := c.Query("token"); len(token) != 0 {
			jwt.GetInstance().MiddlewareFunc()(ctx, c)
		}
	}}, controller.Feed)...)
	{
		_user := root.Group("/user")
		// 用户信息
		_user.GET("/", controller.UserInfo)
		// 登陆, 使用 Hertz 中间价提供的处理方法
		_user.POST("/login/", jwt.GetInstance().LoginHandler)
		// 注册
		_user.POST("/register/", controller.UserRegister)

		_publish := root.Group("/publish", jwt.GetInstance().MiddlewareFunc())
		// 视频投稿
		_publish.POST("/action/", controller.Publish)
		// 获取视频列表
		_publish.GET("/list/", controller.PublishList)

		_favourite := root.Group("/favourite", jwt.GetInstance().MiddlewareFunc())
		// 视频点赞/取消点赞
		_favourite.POST("/action/", UnsupportedMethod)
		// 喜欢视频列表
		_favourite.GET("/list/", UnsupportedMethod)

		_comment := root.Group("/comment", jwt.GetInstance().MiddlewareFunc())
		// 发表评论
		_comment.POST("/action/", UnsupportedMethod)
		// 评论列表
		_comment.GET("/list/", UnsupportedMethod)

		_relation := root.Group("/relation", jwt.GetInstance().MiddlewareFunc())
		// 关注/取消关注
		_relation.POST("/action/", UnsupportedMethod)
		// 关注者列表
		_relation.GET("/follow/list/", UnsupportedMethod)
		// 粉丝列表
		_relation.GET("/follower/list/", UnsupportedMethod)
		// 好友列表
		_relation.GET("/friend/list/", UnsupportedMethod)

		_message := root.Group("/message", jwt.GetInstance().MiddlewareFunc())
		// 轮训获取消息
		_message.GET("/chat/", UnsupportedMethod)
		// 发送消息
		_message.POST("/action/", UnsupportedMethod)
	}
}

func UnsupportedMethod(_ context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, biz.NewFailureResponse("暂不支持该接口服务"))
}
