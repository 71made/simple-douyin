package pack

import (
	"simple-main/http-rcp/cmd/api/biz"
	usvr "simple-main/http-rcp/grpc_gen/user"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/21 20:56
 @ProductName: user.go
 @Description:
*/

func BizBaseUser(u *usvr.User, isFollow bool) *biz.BaseUser {
	if u == nil {
		return nil
	}

	return &biz.BaseUser{
		Id:        u.Id,
		Name:      u.Name,
		AvatarURL: configs.ServerAddr + configs.AvatarURIPrefix + u.Avatar,
		IsFollow:  isFollow,
	}
}

func BizUser(u *usvr.User, isFollow bool) *biz.User {
	if u == nil {
		return nil
	}

	return &biz.User{
		BaseUser:      BizBaseUser(u, isFollow),
		WorkCount:     u.VideoCount,
		LikeCount:     u.FavoriteCount,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
	}
}
