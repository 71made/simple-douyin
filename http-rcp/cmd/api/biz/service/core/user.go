package core

import (
	"context"
	"mime/multipart"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/pack"
	"simple-main/http-rcp/cmd/api/rpc"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/15 11:43
 @ProductName: user.go
 @Description:
*/

type UserLoginRequest struct {
	Username string `query:"username,required" form:"username,required"`
	Password string `query:"password,required" form:"password,required"`
}

type UserLoginResponse struct {
	*biz.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type UserRegisterRequest struct {
	UserLoginRequest
	Avatar *multipart.FileHeader `query:"avatar" form:"avatar"`
}

type UserInfoResponse struct {
	*biz.Response
	User *biz.User `json:"user,omitempty"`
}

// UserService 用户服务接口, 提供 Register Login 和 UserInfo 接口方法
type UserService interface {
	Register(ctx context.Context, req *UserRegisterRequest) (resp *UserLoginResponse)
	Login(ctx context.Context, req *UserLoginRequest) (resp *UserLoginResponse)
	UserInfo(ctx context.Context, userId, thisUserId int64) (resp *UserInfoResponse)
}

func UserServiceImpl() UserService {
	return usInstance
}

// userServiceImpl 对应服务接口实现
type userServiceImpl struct{}

var usInstance = &userServiceImpl{}

// Register 用户注册功能.
// 处理了重复用户创建, 并对用户密码使用 MD5 摘要处理
func (us *userServiceImpl) Register(ctx context.Context, req *UserRegisterRequest) (resp *UserLoginResponse) {
	resp = &UserLoginResponse{}

	// 上传头像
	var avatarURI = configs.EmptyAvatarName
	if req.Avatar != nil {
		// 暂时不做处理, 都使用默认头像
	}
	user, baseResp := rpc.CreateUser(ctx, req.Username, req.Password, avatarURI)

	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp.Response = baseResp
		return
	}

	resp.Response = biz.NewSuccessResponse("注册成功")
	resp.UserId = user.Id
	return
}

// Login 用户登陆功能.
// 使用了基于 MD5 摘要处理的密码判断
// 摘要处理下放至 user server 端服务中
func (us *userServiceImpl) Login(ctx context.Context, req *UserLoginRequest) (resp *UserLoginResponse) {
	resp = &UserLoginResponse{}

	userId, baseResp := rpc.CheckLoginUser(ctx, req.Username, req.Password)
	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp.Response = baseResp
		return
	}

	if userId == biz.NotLoginUserId {
		resp.Response = biz.NewFailureResponse("该用户数据异常")
		return
	}

	resp.Response = biz.NewSuccessResponse("登陆成功")
	resp.UserId = userId
	return
}

func (us *userServiceImpl) UserInfo(ctx context.Context, userId, thisUserId int64) (resp *UserInfoResponse) {
	resp = &UserInfoResponse{}

	user, baseResp := rpc.QueryUser(ctx, userId)
	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp.Response = baseResp
		return
	}

	if user == nil {
		resp.Response = biz.NewFailureResponse("该用户不存在")
		return
	}

	// 查询是否关注
	var isFollow bool
	if thisUserId != biz.NotLoginUserId {
		// 调用 follow server 端服务
	}

	resp.Response = biz.NewSuccessResponse("获取用户信息成功")
	resp.User = pack.BizUser(user, isFollow)
	return
}
