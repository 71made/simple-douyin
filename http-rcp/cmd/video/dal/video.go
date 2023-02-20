package dal

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"simple-main/http-rcp/cmd/video/dal/db"
	"simple-main/http-rcp/pkg/configs"
	"strings"
)

/*
@Author: 71made
@Date: 2023/02/17 22:30
@ProductName: video.go
@Description:
*/

type Video struct {
	gorm.Model
	AuthorId      uint
	PlayUri       string
	CoverUri      string
	FavoriteCount int64
	CommentCount  int64
	Title         string
}

func (v *Video) TableName() string {
	return configs.VideoTable
}

// BeforeCreate
// 通过 GORM 提供的 Hook 实现关联更新 user 记录的 video_count
func (v *Video) BeforeCreate(tx *gorm.DB) (err error) {
	updateRes := tx.Table(configs.UserTable).Where("id = ?", v.AuthorId).
		Update("video_count", gorm.Expr("`video_count` + 1"))

	if updateRes.RowsAffected <= 0 {
		return errors.New("update user record fail")
	}
	if updateRes.RowsAffected > 1 {
		// 做兜底处理
		return errors.New("user table records is dirty")
	}
	return nil
}

func QueryVideos(ctx context.Context, options ...PageOption) ([]*Video, error) {
	res := make([]*Video, 0)

	page := DefaultPage()
	for _, opt := range options {
		opt(page)
	}

	if err := page.Exec(
		db.GetInstance().WithContext(ctx)).
		Order("created_at DESC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryVideosByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("author_id = ?", userId).Order("created_at DESC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryVideosById(ctx context.Context, videoIds []int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if len(videoIds) == 0 {
		return res, nil
	}

	// 构造排序条件
	str := strings.ReplaceAll(fmt.Sprintf("%v", videoIds), " ", ",")
	// 截取中间 id 序列
	str = str[1 : len(str)-1]
	if err := db.GetInstance().WithContext(ctx).
		Where("id in ?", videoIds).Order(fmt.Sprintf("Field(id, %s)", str)).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func IsExistVideo(ctx context.Context, userId int64, title string) (bool, error) {
	ids := make([]int64, 0)
	if err := db.GetInstance().WithContext(ctx).Model(&Video{}).
		Select("id").
		Where("author_id = ? ", userId).Where("title = ?", title).
		Find(&ids).Error; err != nil {
		return false, err
	}
	return len(ids) != 0, nil
}

func CreateVideo(ctx context.Context, video *Video) error {
	return db.GetInstance().WithContext(ctx).Create(video).Error
}
