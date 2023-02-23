package core

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/model"
	"simple-main/simple-http/pkg/configs"
	"strings"
)

/*
 @Author: 71made
 @Date: 2023/01/24 23:10
 @ProductName: user.go
 @Description: 用户相关服务功能
*/

type UserResponse struct {
	biz.Response
	User *biz.User `json:"user,omitempty"`
}

type UserLoginResponse struct {
	biz.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

// UserService 用户服务接口, 提供 Register Login 和 UserInfo 接口方法
type UserService interface {
	Register(ctx context.Context, user *model.User) (resp *UserLoginResponse)
	Login(ctx context.Context, username string, password string) (resp *UserLoginResponse)
	UserInfo(ctx context.Context, userId, thisUserId int64) (resp *UserResponse)
}

func UserServiceImpl() UserService {
	return usInstance
}

// userServiceImpl 对应服务接口实现
type userServiceImpl struct{}

var usInstance = &userServiceImpl{}

// Register 用户注册功能.
// 处理了重复用户创建, 并对用户密码使用 MD5 摘要处理
func (us *userServiceImpl) Register(ctx context.Context, user *model.User) (resp *UserLoginResponse) {
	resp = &UserLoginResponse{}
	h := md5.New()
	if _, err := io.WriteString(h, user.Password); err != nil {
		resp.Response = *biz.NewErrorResponse(err)
		return
	}
	// MD5 摘要算法处理密码
	password := fmt.Sprintf("%x", h.Sum(nil))
	// 回设密码
	user.Password = password

	flag, err := model.IsExistUser(ctx, user.Username)
	if err != nil {
		resp.Response = *biz.NewErrorResponse(err)
		return
	}
	if flag {
		resp.Response = *biz.NewFailureResponse("该用户名已被使用")
		return
	}

	// 设置默认头像
	if len(user.Avatar) == 0 {
		user.Avatar = configs.EmptyAvatarName
	}

	if err := model.CreateUser(ctx, user); err != nil {
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.Response = *biz.NewSuccessResponse("注册成功")
	resp.UserId = int64(user.ID)
	return
}

// Login 用户登陆功能.
// 使用了基于 MD5 摘要处理的密码判断
func (us *userServiceImpl) Login(ctx context.Context, username string, password string) (resp *UserLoginResponse) {
	resp = &UserLoginResponse{}
	h := md5.New()
	if _, err := io.WriteString(h, password); err != nil {
		resp.Response = *biz.NewErrorResponse(err)
		return
	}
	// MD5 摘要算法处理密码
	password = fmt.Sprintf("%x", h.Sum(nil))

	user, err := model.QueryUser(ctx, username)
	if err != nil {
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	if user == nil {
		resp.Response = *biz.NewFailureResponse("该用户暂未注册")
		return
	}

	if strings.Compare(user.Password, password) != 0 {
		resp.Response = *biz.NewFailureResponse("用户名或密码错误")
		return
	}

	resp.Response = *biz.NewSuccessResponse("登陆成功")
	resp.UserId = int64(user.ID)
	return
}

func (us *userServiceImpl) UserInfo(ctx context.Context, userId, thisUserId int64) (resp *UserResponse) {
	resp = &UserResponse{}
	user, err := model.QueryUserById(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	if user == nil {
		resp.Response = *biz.NewFailureResponse("该用户不存在")
		return
	}

	// 查询是否关注
	var isFollow bool
	if thisUserId != biz.NotLoginUserId {
		relation, err := model.QueryRelation(ctx, thisUserId, userId)
		if err != nil {
			hlog.Error(err)
			resp.Response = *biz.NewErrorResponse(err)
			return
		}
		if relation != nil {
			isFollow = relation.IsFollowing()
		}
	}

	resp.Response = *biz.NewSuccessResponse("获取用户信息成功")
	resp.User = &biz.User{
		Id:                 userId,
		Name:               user.Username,
		AvatarURL:          configs.ServerAddr + configs.AvatarURIPrefix + user.Avatar,
		FollowCount:        user.FollowCount,
		FollowerCount:      user.FollowerCount,
		WorkCount:          user.VideoCount,
		BackgroundImage:    user.BackgroundImage,
		Signature:          user.Signature,
		FavoriteCount:      user.FavoriteCount,
		TotalFavoriteCount: user.TotalFavoriteCount,
		IsFollow:           isFollow,
	}
	return
}
