package pack

import (
	"simple-main/http-rcp/cmd/favorite/dal"
	fsvr "simple-main/http-rcp/grpc_gen/favorite"
)

/*
 @Author: 71made
 @Date: 2023/02/21 02:51
 @ProductName: favorite.go
 @Description:
*/

func Favorite(f *dal.Favorite) *fsvr.Favorite {
	if f == nil {
		return nil
	}

	return &fsvr.Favorite{
		Id:         int64(f.ID),
		UserId:     int64(f.UserId),
		VideoId:    int64(f.VideoId),
		IsFavorite: f.IsFavorite(),
	}
}

func Favorites(fs []*dal.Favorite) []*fsvr.Favorite {
	if fs == nil || len(fs) == 0 {
		return make([]*fsvr.Favorite, 0)
	}

	res := make([]*fsvr.Favorite, len(fs))
	for i, f := range fs {
		res[i] = Favorite(f)
	}

	return res
}
