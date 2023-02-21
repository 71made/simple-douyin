package first

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/pack"
	"simple-main/http-rcp/cmd/api/rpc"
)

/*
 @Author: 71made
 @Date: 2023/02/21 14:36
 @ProductName: favorite.go
 @Description:
*/

type FavoriteActionRequest struct {
	UserId     int64 `query:"user_id"`
	VideoId    int64 `query:"video_id,required"`
	ActionType uint  `query:"action_type,required"`
}

type FavoriteListResponse struct {
	*biz.Response
	VideoList []*biz.Video `json:"video_list"`
}

type FavoriteService interface {
	Action(ctx context.Context, req *FavoriteActionRequest) (resp *biz.Response)
	FavoriteList(ctx context.Context, userId int64) (resp *FavoriteListResponse)
}

func FavoriteServiceImpl() FavoriteService {
	return fsInstance
}

// favoriteServiceImpl 对应服务接口实现
type favoriteServiceImpl struct{}

var fsInstance = &favoriteServiceImpl{}

func (fs *favoriteServiceImpl) Action(ctx context.Context, req *FavoriteActionRequest) (resp *biz.Response) {

	baseResp := rpc.FavoriteAction(ctx, req.UserId, req.VideoId, req.ActionType)

	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp = biz.NewFailureResponse(baseResp.StatusMsg)
		return
	}

	resp = biz.NewSuccessResponse("操作成功")
	return
}

func (fs *favoriteServiceImpl) FavoriteList(ctx context.Context, userId int64) (resp *FavoriteListResponse) {
	resp = &FavoriteListResponse{}

	videos, baseResp := rpc.QueryFavoriteVideos(ctx, userId)
	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp.VideoList = make([]*biz.Video, 0)
		resp.Response = baseResp
		return
	}

	if len(videos) == 0 {
		resp.VideoList = make([]*biz.Video, 0)
		resp.Response = biz.NewSuccessResponse("获取成功")
		return
	}

	// 转换为 []biz.Video
	videoList, err := pack.BizVideos(ctx, videos, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = biz.NewErrorResponse(err)
		return
	}

	resp.VideoList = videoList
	resp.Response = biz.NewSuccessResponse("获取成功")
	return
}
