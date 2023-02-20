package core

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/pack"
	"simple-main/http-rcp/cmd/api/rpc"
	"simple-main/http-rcp/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/18 22:40
 @ProductName: publish.go
 @Description:
*/

type VideoListResponse struct {
	*biz.Response
	VideoList []*biz.Video `json:"video_list"`
}

type VideoPublishRequest struct {
	VideoFinalName string
	UserId         int64
	Title          string
}

// PublishService 视频服务接口, 提供 Register Login 和 UserInfo 接口方法
type PublishService interface {
	GetPublishList(ctx context.Context, userId int64) (resp *VideoListResponse)
	PublishVideo(ctx context.Context, req *VideoPublishRequest) (resp *biz.Response)
}

func PublishServiceImpl() PublishService {
	return psInstance
}

type publishServiceImpl struct{}

var psInstance = &publishServiceImpl{}

func (ps *publishServiceImpl) GetPublishList(ctx context.Context, userId int64) (resp *VideoListResponse) {

	resp = &VideoListResponse{}

	videos, baseResp := rpc.QueryVideosByUserId(ctx, userId)
	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp.VideoList = make([]*biz.Video, 0)
		resp.Response = baseResp
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

func (ps *publishServiceImpl) PublishVideo(ctx context.Context, req *VideoPublishRequest) (resp *biz.Response) {
	defer removeCache(configs.VideoPathPrefix + req.VideoFinalName)
	return rpc.CreateVideo(ctx, req.VideoFinalName, req.UserId, req.Title)
}

// removeCache 删除缓存资源
func removeCache(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		_ = os.Remove(filePath)
	}
}
