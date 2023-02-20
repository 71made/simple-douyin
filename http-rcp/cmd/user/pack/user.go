package pack

import (
	"simple-main/http-rcp/cmd/user/dal"
	usvr "simple-main/http-rcp/grpc_gen/user"
)

/*
 @Author: 71made
 @Date: 2023/02/15 12:18
 @ProductName: user.go
 @Description:
*/

func User(u *dal.User) *usvr.User {
	if u == nil {
		return nil
	}

	return &usvr.User{
		Id:                 int64(u.ID),
		Name:               u.Username,
		Avatar:             u.Avatar,
		Signature:          u.Signature,
		BackgroundImage:    u.BackgroundImage,
		VideoCount:         u.VideoCount,
		FavoriteCount:      u.FavoriteCount,
		TotalFavoriteCount: u.TotalFavoriteCount,
		FollowCount:        u.FollowCount,
		FollowerCount:      u.FollowerCount,
	}
}

func Users(us []*dal.User) []*usvr.User {
	if us == nil || len(us) == 0 {
		return make([]*usvr.User, 0)
	}
	res := make([]*usvr.User, len(us))
	for i, u := range us {
		res[i] = User(u)
	}

	return res
}
