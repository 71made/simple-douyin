package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"simple-main/http-rcp/cmd/video/dal"
	"simple-main/http-rcp/cmd/video/pack"
	vsvr "simple-main/http-rcp/grpc_gen/video"
	"simple-main/http-rcp/pkg/configs"
	"simple-main/http-rcp/pkg/utils/minio"
	"strings"
)

/*
 @Author: 71made
 @Date: 2023/02/17 22:30
 @ProductName: handler.go
 @Description:
*/

func newServer(opts ...grpc.ServerOption) *grpc.Server {
	svr := grpc.NewServer(opts...)
	vsvr.RegisterVideoManagementServer(svr, &VideoManagementServer{})
	return svr
}

type VideoManagementServer struct {
	vsvr.UnimplementedVideoManagementServer
}

func (vss VideoManagementServer) CreateVideo(ctx context.Context, req *vsvr.CreateVideoRequest) (*vsvr.CreateVideoResponse, error) {
	resp := &vsvr.CreateVideoResponse{}
	found, _ := dal.IsExistVideo(ctx, req.AuthorId, req.Title)
	if found {
		resp.BaseResponse = pack.NewFailureResponse("已存在同标题视频投稿")
		return resp, nil
	}

	// 并发上传视频和封面
	// 通过 channel 并发控制
	videoUriChan := make(chan string)
	coverUriChan := make(chan string)
	var uploadVideoErr, uploadCoverErr error
	go func() {
		err := minio.UploadVideo(ctx, req.VideoFinalName, filepath.Join(configs.VideoPathPrefix, req.VideoFinalName))
		if err != nil {
			grpclog.Error(err)
			uploadVideoErr = err
			videoUriChan <- ""
		}
		videoUriChan <- req.VideoFinalName
	}()

	go func() {
		coverFinalName, err := readFrameAsJpeg(req.VideoFinalName)
		if err != nil {
			grpclog.Error(err)
			uploadCoverErr = err
			coverUriChan <- configs.EmptyCoverName
		}
		err = minio.UploadCover(ctx, coverFinalName, filepath.Join(configs.CoverPathPrefix, coverFinalName))
		if err != nil {
			grpclog.Error(err)
			coverUriChan <- configs.EmptyCoverName
		}
		coverUriChan <- coverFinalName
	}()

	// 处理异常错误
	var err error
	if uploadVideoErr != nil {
		err = uploadVideoErr
	}
	if uploadCoverErr != nil {
		err = uploadCoverErr
	}

	videoUri := <-videoUriChan
	coverUri := <-coverUriChan

	// 删除服务器缓存
	defer func() {
		if coverUri != configs.EmptyCoverName {
			removeCache(filepath.Join(configs.CoverPathPrefix + coverUri))
		}
	}()

	defer func() {
		if err != nil {
			if len(videoUri) != 0 {
				go func() {
					_ = minio.RemoveVideo(ctx, videoUri)
				}()
			}
			if coverUri != configs.EmptyCoverName {
				go func() {
					_ = minio.RemoveCover(ctx, coverUri)
				}()
			}
		}
	}()

	if len(videoUri) == 0 {
		resp.BaseResponse = pack.NewFailureResponse("视频上传出错, 请稍后重试")
		err = errors.New("upload video error")
		return resp, nil
	}

	// 构建实体
	video := &dal.Video{
		AuthorId: uint(req.AuthorId),
		PlayUri:  videoUri,
		CoverUri: coverUri,
		Title:    req.Title,
	}

	if err = dal.CreateVideo(ctx, video); err != nil {
		resp.BaseResponse = pack.NewErrorResponse(err)
		return resp, nil
	}

	resp.BaseResponse = pack.NewSuccessResponse("投稿成功")
	return resp, nil
}

func (vss VideoManagementServer) QueryVideos(ctx context.Context, req *vsvr.QueryVideosRequest) (*vsvr.QueryVideosResponse, error) {

	// 内部抽取处理函数
	packQueryRes := func(videos []*dal.Video, err error) *vsvr.QueryVideosResponse {
		resp := &vsvr.QueryVideosResponse{}
		if err != nil {
			grpclog.Error(err)
			resp.BaseResponse = pack.NewErrorResponse(err)
			resp.VideoList = make([]*vsvr.Video, 0)
			return resp
		}

		resp.VideoList = pack.Videos(videos)
		resp.BaseResponse = pack.NewSuccessResponse("获取成功")
		return resp
	}

	if req.UserId != nil {
		videos, err := dal.QueryVideosByUserId(ctx, *req.UserId)
		return packQueryRes(videos, err), nil
	}

	if req.VideoIds != nil && len(req.VideoIds) != 0 {
		videos, err := dal.QueryVideosById(ctx, req.VideoIds)
		return packQueryRes(videos, err), nil
	}

	return packQueryRes(make([]*dal.Video, 0), errors.New("请求服务参数异常缺失")), nil
}

func (vss VideoManagementServer) QueryFeedVideos(ctx context.Context, req *vsvr.QueryFeedVideoRequest) (*vsvr.QueryFeedVideosResponse, error) {
	resp := &vsvr.QueryFeedVideosResponse{}

	videos, err := dal.QueryVideos(ctx, dal.PageAfter(req.LastTime), dal.PageLimit(int(req.Limit)))
	if err != nil {
		grpclog.Error(err)
		resp.BaseResponse = pack.NewErrorResponse(err)
		resp.VideoList = make([]*vsvr.Video, 0)
		return resp, nil
	}

	resp.VideoList = pack.Videos(videos)
	resp.BaseResponse = pack.NewSuccessResponse("获取成功")
	return resp, nil
}

// readFrameAsJpeg 截取视频封面
func readFrameAsJpeg(videoFinalName string) (coverFinalName string, err error) {
	videoPath := configs.VideoPathPrefix + videoFinalName

	coverFinalName = strings.TrimSuffix(videoFinalName, path.Ext(videoFinalName)) + ".jpeg"
	coverPath := configs.CoverPathPrefix + coverFinalName

	// 使用 ffmpeg 提取指定帧作为图像文件
	cmd := exec.Command("ffmpeg",
		"-y", // 强制覆盖
		"-i",
		videoPath,                 // 视频路径
		"-vf", "select=eq(n\\,1)", // 抽取第 2 帧
		"-vframes", "1",
		coverPath, // 输出封面路径
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("error output: %s", string(output))
	}
	return
}

// removeCache 删除缓存资源
func removeCache(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		_ = os.Remove(filePath)
	}
}
