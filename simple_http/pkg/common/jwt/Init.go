package jwt

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"net/http"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/biz/controller/core"
	service "simple-main/simple-http/cmd/biz/service/core"
	"simple-main/simple-http/pkg/configs"
	"sync"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/01/30 21:44
 @ProductName: Init.go
 @Description: 初始化 Hertz 中间件 JWT
*/

var once sync.Once
var jwtMiddleware *jwt.HertzJWTMiddleware

func Init() {
	jwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:           []byte(configs.JWTSecretKey),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt, form: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Timeout:       24 * time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   configs.IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &biz.User{
				Id: int64(claims[configs.IdentityKey].(float64)),
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					configs.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			core.UserLogin(ctx, c)
			resp, found := c.Get("resp")
			loginResp := resp.(*service.UserLoginResponse)
			if found && loginResp.StatusCode == 0 {
				return loginResp.UserId, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			status, _ := c.Get("status")
			if status == nil {
				status = http.StatusOK
			}

			resp, found := c.Get("resp")
			loginResp := resp.(*service.UserLoginResponse)
			if found && loginResp.StatusCode == 0 {
				loginResp.Token = token
			}
			c.JSON(status.(int), loginResp)
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			status, _ := c.Get("status")
			if status == nil {
				status = code
			}
			resp, found := c.Get("resp")

			if found {
				c.JSON(status.(int), resp)
			} else {
				c.JSON(status.(int), biz.NewFailureResponse(message))
			}
		},
	})
}

func GetInstance() *jwt.HertzJWTMiddleware {

	if jwtMiddleware == nil {
		once.Do(Init)
	}
	return jwtMiddleware
}
