package first

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/cmd/biz"
	"simple-main/cmd/configs"
	"simple-main/cmd/model"
)

/*
 @Author: 71made
 @Date: 2023/01/31 15:37
 @ProductName: favorite.go
 @Description:
*/

type FavoriteRequest struct {
	UserId     int64 `json:"user_id"`
	VideoId    int64 `json:"video_id"`
	ActionType uint  `json:"action_type"`
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
	} else if found.FavoriteType != f.FavoriteType {
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
		}

		// 喜欢列表都为 true
		videoList[i].IsFavorite = true
		videoIds[i] = int64(video.ID)
		videoMap[video.ID] = &videoList[i]
	}
	// 填充剩下参数
	// 查询是否关注
	//author.IsFollow = false
	return videoList, nil
}
