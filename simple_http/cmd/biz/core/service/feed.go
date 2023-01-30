package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/cmd/biz"
	"simple-main/cmd/model"
	"simple-main/pkg/configs"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/01/28 17:21
 @ProductName: feed.go
 @Description:
*/

type FeedResponse struct {
	biz.Response
	VideoList []biz.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

// FeedService 视频流服务接口
type FeedService interface {
	GetFeed(ctx context.Context, lastTime time.Time, userId int64) (resp *FeedResponse)
}

func GetFeedServiceImpl() FeedService {
	return fsInstance
}

type feedServiceImpl struct{}

var fsInstance = &feedServiceImpl{}

func (fs *feedServiceImpl) GetFeed(ctx context.Context, lastTime time.Time, userId int64) (resp *FeedResponse) {
	resp = &FeedResponse{}

	videos, err := model.QueryVideos(ctx,
		model.PageLimit(30),
		model.PageAfter(lastTime),
	)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	// 转换为 []biz.Video
	videoList, err := GetBizVideoList(ctx, videos)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.VideoList = videoList
	resp.NextTime = videos[len(videos)-1].CreatedAt.Unix()
	resp.Response = *biz.NewSuccessResponse("获取成功")
	return
}

func GetBizVideoList(ctx context.Context, videos []model.Video) ([]biz.Video, error) {
	// 构建转换实体
	var videoList = make([]biz.Video, len(videos))
	// 缓存 author
	var authors = make(map[uint]biz.User)
	for i, video := range videos {
		author, found := authors[video.AuthorId]
		if !found {
			user, err := model.QueryUserById(ctx, int64(video.AuthorId))
			if err != nil {
				return nil, err
			}
			author = biz.User{
				Id:            int64(video.AuthorId),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				//IsFollow:      false,
			}
			authors[video.AuthorId] = author
		}

		videoList[i] = biz.Video{
			Id:            int64(video.ID),
			Author:        author,
			PlayUrl:       configs.ServerAddr + configs.VideoUriPrefix + video.PlayUri,
			CoverUrl:      configs.ServerAddr + configs.CoverUriPrefix + video.CoverUri,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			//IsFavorite:    false,
		}
	}
	return videoList, nil
}
