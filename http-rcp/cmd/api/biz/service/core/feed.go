package core

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/pack"
	"simple-main/http-rcp/cmd/api/rpc"
)

/*
 @Author: 71made
 @Date: 2023/02/19 02:18
 @ProductName: feed.go
 @Description:
*/

type FeedResponse struct {
	*biz.Response
	VideoList []*biz.Video `json:"video_list,omitempty"`
	NextTime  int64        `json:"next_time,omitempty"`
}

// FeedService 视频流服务接口
type FeedService interface {
	GetFeed(ctx context.Context, lastTime int64, userId int64) (resp *FeedResponse)
}

func FeedServiceImpl() FeedService {
	return fsInstance
}

type feedServiceImpl struct{}

var fsInstance = &feedServiceImpl{}

func (fs *feedServiceImpl) GetFeed(ctx context.Context, lastTime int64, userId int64) (resp *FeedResponse) {
	resp = &FeedResponse{}

	videos, baseResp := rpc.QueryFeedVideos(ctx, 30, lastTime)
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

	// 转换为 []*biz.Video
	videoList, err := pack.BizVideos(ctx, videos, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = biz.NewErrorResponse(err)
		return
	}

	resp.VideoList = videoList
	resp.NextTime = videos[len(videos)-1].CreateTime
	resp.Response = biz.NewSuccessResponse("获取成功")
	return
}
