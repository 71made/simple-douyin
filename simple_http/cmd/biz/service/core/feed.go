package core

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/configs"
	"simple-main/simple-http/cmd/model"
	"sync"
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

// GetBizVideoList
// 类型转换 []model.Video -> []biz.Video, 用户未登陆时, userId 传递 NotLoginUserId
func GetBizVideoList(ctx context.Context, videos []model.Video, userId int64) ([]biz.Video, error) {

	// 缓存 author 映射关系
	var authors = make(map[uint]*biz.User)
	// 缓存 author id
	var authorIds = make([]int64, len(authors))

	// 缓存 video id 为后续查询关系做准备
	var videoIds = make([]int64, len(videos))
	// 保存 video id 与实体间的映射关系
	var videoMap = make(map[uint]*biz.Video, len(videos))

	// 保存关注情况
	var isFollowMap = make(map[uint]bool)

	// 构建转换实体
	var videoList = make([]biz.Video, len(videos))
	for i, video := range videos {
		// 保存 authorIds 和初始化 authors、isFollowMap
		_, found := authors[video.AuthorId]
		if !found {
			authors[video.AuthorId] = nil
			authorIds = append(authorIds, int64(video.AuthorId))
			isFollowMap[video.AuthorId] = false
		}

		// 构造 video
		videoList[i] = biz.Video{
			Id:            int64(video.ID),
			PlayURL:       configs.ServerAddr + configs.VideoURIPrefix + video.PlayUri,
			CoverURL:      configs.ServerAddr + configs.CoverURIPrefix + video.CoverUri,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
		}
		videoIds[i] = int64(video.ID)
		videoMap[video.ID] = &videoList[i]
	}

	// 并发控制
	var wg sync.WaitGroup
	wg.Add(1) // 后续查询 authors 实体的并发任务

	// 对于已登陆用户, 需要查询用户对应的关注和点赞情况
	if userId != biz.NotLoginUserId {
		wg.Add(2) // 接下来查询关注和点赞的并发任务

		go func() {
			// 查询是否关注
			relations, _ := model.QueryRelations(ctx, userId, authorIds)

			for _, relation := range relations {
				isFollowMap[relation.UserId] = relation.IsFollowing()
			}
			wg.Done()
		}()

		go func() {
			// 查询是否点赞
			favorites, _ := model.QueryFavorites(ctx, userId, videoIds)

			for _, favorite := range favorites {
				video := videoMap[favorite.VideoId]
				video.IsFavorite = favorite.IsFavorite()
			}
			wg.Done()
		}()

	}

	var QUserErr error
	go func() {
		defer wg.Done()
		// 批量查询
		users, err := model.QueryUsersByIds(ctx, authorIds)
		if QUserErr = err; err != nil {
			return
		}
		for _, user := range users {
			author := &biz.User{
				Id:            int64(user.ID),
				Name:          user.Username,
				AvatarURL:     configs.ServerAddr + configs.AvatarURIPrefix + user.Avatar,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
			}
			authors[user.ID] = author
		}
	}()

	wg.Wait()

	// 处理异常错误
	if QUserErr != nil {
		return nil, QUserErr
	}

	// 最后组装
	for i, video := range videos {
		author := authors[video.AuthorId]
		author.IsFollow = isFollowMap[video.AuthorId]
		videoList[i].Author = *author
	}

	return videoList, nil
}
