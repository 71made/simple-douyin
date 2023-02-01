package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"simple-main/cmd/biz"
	"simple-main/cmd/common/minio"
	"simple-main/cmd/configs"
	"simple-main/cmd/model"
	"strings"
)

/*
 @Author: 71made
 @Date: 2023/01/28 16:30
 @ProductName: publish.go
 @Description:
*/

type VideoListResponse struct {
	biz.Response
	VideoList []biz.Video `json:"video_list"`
}

type VideoPublishRequest struct {
	VideoFinalName string
	UserId         int64
	Title          string
}

// PublishService 视频服务接口, 提供 Register Login 和 UserInfo 接口方法
type PublishService interface {
	GetPublishList(ctx context.Context, userId int64) (resp *VideoListResponse)
	PublishVideo(ctx context.Context, req *VideoPublishRequest, method UploadMethod) (resp *biz.Response)
}

func GetPublishServiceImpl() PublishService {
	return psInstance
}

type publishServiceImpl struct{}

var psInstance = &publishServiceImpl{}

func (ps *publishServiceImpl) GetPublishList(ctx context.Context, userId int64) (resp *VideoListResponse) {

	resp = &VideoListResponse{}

	videos, err := model.QueryVideosByUserId(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
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
	resp.Response = *biz.NewSuccessResponse("获取成功")
	return
}

func (ps *publishServiceImpl) PublishVideo(ctx context.Context, req *VideoPublishRequest, Upload UploadMethod) (resp *biz.Response) {

	// 检查是否存在同标题视频
	flag, err := model.IsExistVideo(ctx, req.UserId, req.Title)
	if err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(err)
		return
	}
	if flag {
		resp = biz.NewFailureResponse("已存在同标题视频投稿")
		return
	}

	// 构建视频保存服务器路径
	videoFilePath := filepath.Join(configs.VideoPathPrefix, req.VideoFinalName)
	// 执行视频保存服务器和上传 MinIO
	if err := Upload(videoFilePath, func(data *multipart.FileHeader) error {
		return minio.UploadVideo(ctx, req.VideoFinalName, data)
	}); err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(fmt.Errorf("视频上传出错"))
		return
	}

	// 获取并上传视频封面至 MinIO
	nameChan := getAndUploadCover(ctx, req)

	// 构建 video 记录
	newVideo := model.Video{
		AuthorId: uint(req.UserId),
		PlayUri:  req.VideoFinalName,
		CoverUri: configs.EmptyCoverName, // 为保证都投稿视频能有封面, 先设置为默认封面, 封面上传成功后再更新
		Title:    req.Title,
	}

	// 创建 video 记录
	if err := model.CreateVideo(ctx, &newVideo); err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(err)
		return
	}

	coverFinalName := <-nameChan
	if coverFinalName != "" {
		// 上传成功, 设置封面
		newVideo.CoverUri = coverFinalName
		err := model.UpdateVideoCoverUri(ctx, &newVideo)
		// 更新失败, 恢复默认封面
		if err != nil {
			newVideo.CoverUri = configs.EmptyCoverName
			err = errors.New("video's cover image update fail")
		}

	}
	hlog.Info(fmt.Sprintf("video path: %s, cover path: %s",
		configs.ServerAddr+configs.VideoUriPrefix+newVideo.PlayUri, configs.ServerAddr+configs.CoverUriPrefix+newVideo.CoverUri))
	resp = biz.NewSuccessResponse("投稿成功")

	// 删除服务器缓存
	defer func() {
		removeCache(videoFilePath)
		if coverFinalName != "" {
			removeCache(filepath.Join(configs.CoverPathPrefix + coverFinalName))
		}
	}()

	// 处理错误
	defer func() {
		if err != nil {
			// 发生错误, 删除上传文件
			if err.Error() != "video's cover image update fail" {
				go func() {
					_ = minio.RemoveVideo(ctx, req.VideoFinalName)
				}()
			}

			if coverFinalName != "" {
				go func() {
					_ = minio.RemoveCover(ctx, coverFinalName)
				}()
			}
		}
	}()
	return
}

// UploadMethod 视频文件上传方法 dstPath 为服务器保存路径, uploadToMinIO 中调用 minio.UploadVideo
type UploadMethod func(dstPath string, uploadToMinIO func(data *multipart.FileHeader) error) error

func getAndUploadCover(ctx context.Context, req *VideoPublishRequest) <-chan string {
	coverFinalName, readErr := readFrameAsJpeg(req.VideoFinalName)
	name := make(chan string)
	// 封面截取出错不影响后续操作
	if readErr != nil {
		hlog.Error(fmt.Errorf("msg: %s, err: %v", "视频封面截取异常", readErr))
		name <- ""
	} else {
		// 上传 MinIO
		go func() {
			uploadErr := minio.UploadCoverWithFilePath(ctx, coverFinalName, filepath.Join(configs.CoverPathPrefix+coverFinalName))
			if uploadErr != nil {
				hlog.Error(fmt.Errorf("msg: %s, err: %v", "视频封面上传出错", uploadErr))
			}
			name <- coverFinalName
		}()

	}
	return name
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

func removeCache(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		_ = os.Remove(filePath)
	}
}
