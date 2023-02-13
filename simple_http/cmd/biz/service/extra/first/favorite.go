package first

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/model"
	"simple-main/simple-http/pkg/configs"
	"sync"
)

/*
 @Author: 71made
 @Date: 2023/01/31 15:37
 @ProductName: favorite.go
 @Description:
*/

type FavoriteRequest struct {
	UserId     int64 `query:"user_id"`
	VideoId    int64 `query:"video_id,required"`
	ActionType uint  `query:"action_type,required"`
}

type FavoriteListResponse struct {
	biz.Response
	VideoList []biz.Video `json:"video_list"`
}

type FavoriteService interface {
	Action(ctx context.Context, req *FavoriteRequest) (resp *biz.Response)
	FavoriteList(ctx context.Context, userId int64) (resp *FavoriteListResponse)
}

func FavoriteServiceImpl() FavoriteService {
	return fsInstance
}

// favoriteServiceImpl 对应服务接口实现
type favoriteServiceImpl struct{}

var fsInstance = &favoriteServiceImpl{}

func (fs *favoriteServiceImpl) Action(ctx context.Context, req *FavoriteRequest) (resp *biz.Response) {
	// 构建实体
	f := &model.Favorite{
		UserId:       uint(req.UserId),
		VideoId:      uint(req.VideoId),
		FavoriteType: req.ActionType,
	}

	// 查找点赞记录
	found, err := model.QueryFavorite(ctx, req.UserId, req.VideoId)
	if err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(err)
		return
	}

	// 没有记录则创建, 有则更新
	if found == nil {
		err = model.CreateFavorite(ctx, f)
	} else if found.FavoriteType != f.GetFavoriteType() {
		// 并且只对于不同的 type, 才触发更新
		err = model.UpdateFavorite(ctx, f)
	}

	if err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(err)
		return
	}
	resp = biz.NewSuccessResponse("操作成功")
	return
}

func (fs *favoriteServiceImpl) FavoriteList(ctx context.Context, userId int64) (resp *FavoriteListResponse) {
	resp = &FavoriteListResponse{}

	videos, err := model.QueryFavoriteVideos(ctx, userId)
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
	videoList, err := GetBizFavoriteVideoList(ctx, videos, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.VideoList = videoList
	resp.Response = *biz.NewSuccessResponse("获取成功")
	return
}

func GetBizFavoriteVideoList(ctx context.Context, videos []model.Video, userId int64) ([]biz.Video, error) {

	// 缓存 author
	var authors = make(map[uint]*biz.User)
	// 保存 author id
	var authorIds = make([]int64, 0)

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

		videoList[i] = biz.Video{
			Id:            int64(video.ID),
			PlayURL:       configs.ServerAddr + configs.VideoURIPrefix + video.PlayUri,
			CoverURL:      configs.ServerAddr + configs.CoverURIPrefix + video.CoverUri,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
		}

		// 喜欢列表都为 true
		videoList[i].IsFavorite = true
		videoMap[video.ID] = &videoList[i]
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// 填充剩下参数
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

	go func() {
		// 查询是否关注
		relations, _ := model.QueryRelations(ctx, userId, authorIds)

		for _, relation := range relations {
			isFollowMap[relation.UserId] = relation.IsFollowing()
		}
		wg.Done()
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
