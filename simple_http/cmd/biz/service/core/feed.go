package core

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/cmd/biz"
	"simple-main/cmd/configs"
	"simple-main/cmd/model"
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

	if len(videos) == 0 {
		resp.VideoList = make([]biz.Video, 0)
		resp.Response = *biz.NewSuccessResponse("获取成功")
		return
	}

	// 转换为 []biz.Video
	videoList, err := GetBizVideoList(ctx, videos, userId)
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

// NotLoginUserId 定义未登陆用户 id 为 -1
const NotLoginUserId = -1

// GetBizVideoList
// 类型转换 []model.Video -> []biz.Video, 用户未登陆时, userId 传递 NotLoginUserId
func GetBizVideoList(ctx context.Context, videos []model.Video, userId int64) ([]biz.Video, error) {
	// 构建转换实体
	var videoList = make([]biz.Video, len(videos))
	// 缓存 author
	var authors = make(map[uint]biz.User)
	// 缓存 video id 为后续查询关系做准备
	var videoIds = make([]int64, len(videos))
	// 保存 video id 与实体间的映射关系
	var videoMap = make(map[uint]*biz.Video, len(videos))
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
		videoIds[i] = int64(video.ID)
		videoMap[video.ID] = &videoList[i]
	}
	// 填充剩下参数
	if userId != NotLoginUserId {
		// 查询是否关注
		//author.IsFollow = false
		// 查询是否点赞
		favorites, _ := model.QueryFavorites(ctx, userId, videoIds)

		for _, favorite := range favorites {
			video := videoMap[favorite.VideoId]
			video.IsFavorite = favorite.IsFavorite()
		}

	}
	return videoList, nil
}
