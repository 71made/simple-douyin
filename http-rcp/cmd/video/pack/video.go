package pack

import (
	"simple-main/http-rcp/cmd/video/dal"
	vsvr "simple-main/http-rcp/grpc_gen/video"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/17 22:32
 @ProductName: video.go
 @Description:
*/

func Video(v *dal.Video) *vsvr.Video {
	if v == nil {
		return nil
	}

	return &vsvr.Video{
		Id:            int64(v.ID),
		AuthorId:      int64(v.AuthorId),
		PlayUrl:       configs.ServerAddr + configs.VideoURIPrefix + v.PlayUri,
		CoverUrl:      configs.ServerAddr + configs.CoverURIPrefix + v.CoverUri,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		Title:         v.Title,
		CreateTime:    v.CreatedAt.Unix(),
	}
}

func Videos(vs []*dal.Video) []*vsvr.Video {
	if vs == nil || len(vs) == 0 {
		return make([]*vsvr.Video, 0)
	}
	res := make([]*vsvr.Video, len(vs))
	for i, v := range vs {
		res[i] = Video(v)
	}

	return res
}
