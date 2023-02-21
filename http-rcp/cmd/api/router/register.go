package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"net/http"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/handler/core"
	"simple-main/http-rcp/cmd/api/biz/handler/extra/first"
	"simple-main/http-rcp/cmd/api/router/jwt"
	"strings"
)

/*
 @Author: 71made
 @Date: 2023/01/25 12:08
 @ProductName: register.go
 @Description:
*/

// register 路由注册. 全部的路由注册配置都在此函数中
func register(r *server.Hertz) {

	// 静态资源
	//r.Static("/static", "./static")

	root := r.Group("/douyin")
	// 获取视频流
	root.GET("/feed/", append([]app.HandlerFunc{func(ctx context.Context, c *app.RequestContext) {
		// 对于 Feed 接口, 如果传入了 token, 此处需要手动调用 JWT 的 mw 校验和解析 token
		if token := c.Query("token"); len(token) != 0 {
			jwt.GetInstance().MiddlewareFunc()(ctx, c)
		}
	}}, core.Feed)...)
	{
		_user := root.Group("/user")
		// 用户信息
		_user.GET("/", core.UserInfo)
		// 登陆, 使用 Hertz 中间价提供的处理方法
		_user.POST("/login/", jwt.GetInstance().LoginHandler)
		// 注册
		_user.POST("/register/", core.UserRegister)

		_publish := root.Group("/publish", jwt.GetInstance().MiddlewareFunc())
		// 视频投稿
		_publish.POST("/action/", core.Publish)
		// 获取视频列表
		_publish.GET("/list/", core.PublishList)

		_favorite := root.Group("/favorite", func(ctx context.Context, c *app.RequestContext) {
			// 对于 /list/  接口, 在用户未登陆时也可以请求查看其他用户的喜欢视频列表
			// 所以如果传入了 token, 此处需要手动调用 JWT 的 mw 校验和解析 token
			// 但对于 /action/ 是需要校验和解析 token 的
			reqURI := c.GetRequest().URI().String()
			if token := c.Query("token"); len(token) != 0 ||
				strings.Contains(reqURI, "/action/") {
				jwt.GetInstance().MiddlewareFunc()(ctx, c)
			}
		})
		// 视频点赞/取消点赞
		_favorite.POST("/action/", first.FavoriteAction)
		// 喜欢视频列表
		_favorite.GET("/list/", first.GetFavoriteList)

		_comment := root.Group("/comment", func(ctx context.Context, c *app.RequestContext) {
			// 对于 /list/  接口, 在用户未登陆时也可以请求查看视频的评论列表
			// 所以如果传入了 token, 此处需要手动调用 JWT 的 mw 校验和解析 token
			// 但对于 /action/ 是需要校验和解析 token 的
			reqURI := c.GetRequest().URI().String()
			if token := c.Query("token"); len(token) != 0 ||
				strings.Contains(reqURI, "/action/") {
				jwt.GetInstance().MiddlewareFunc()(ctx, c)
			}
		})
		// 发表评论
		_comment.POST("/action/", UnsupportedMethod)
		// 评论列表
		_comment.GET("/list/", UnsupportedMethod)

		_relation := root.Group("/relation", func(ctx context.Context, c *app.RequestContext) {
			// 对于 /follow/list/ 和 /follower/list/ 接口, 在用户未登陆时也可以请求查看其他用户的关注和粉丝列表
			// 所以如果传入了 token, 此处需要手动调用 JWT 的 mw 校验和解析 token
			// 但对于 /action/ 和 /friend/list/ 接口, 也是需要校验和解析 token 的
			reqURI := c.GetRequest().URI().String()
			if token := c.Query("token"); len(token) != 0 ||
				strings.Contains(reqURI, "/action/") ||
				strings.Contains(reqURI, "/friend/list/") {
				jwt.GetInstance().MiddlewareFunc()(ctx, c)
			}
		})
		// 关注/取消关注
		_relation.POST("/action/", UnsupportedMethod)
		// 关注者列表
		_relation.GET("/follow/list/", UnsupportedMethod)
		// 粉丝列表
		_relation.GET("/follower/list/", UnsupportedMethod)
		// 好友列表
		_relation.GET("/friend/list/", UnsupportedMethod)

		_message := root.Group("/message", func(ctx context.Context, c *app.RequestContext) {
			// 对于 /chat/ 接口, App 端是不断轮询的;
			// 而在用户未登陆时发送的请求, 当进行 token 校验和解析时
			// 会触发 JWT 中间件的 Unauthorized 回调方法, 返回失败响应,
			// 进而导致 App 端提示 "网络异常，请检查！"
			// 所以此处在未登陆时, 放行 /chat/ 请求, 即不做 token 校验和解析, 交由 service 层处理
			reqURI := c.GetRequest().URI().String()
			if token := c.Query("token"); len(token) != 0 ||
				strings.Contains(reqURI, "/action/") {
				jwt.GetInstance().MiddlewareFunc()(ctx, c)
			}
		})
		// 轮训获取消息
		_message.GET("/chat/", UnsupportedMethod)
		// 发送消息
		_message.POST("/action/", UnsupportedMethod)
	}
}

func UnsupportedMethod(_ context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, biz.NewFailureResponse("暂不支持该接口服务"))
}
