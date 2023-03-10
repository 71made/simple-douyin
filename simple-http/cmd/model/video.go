package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"simple-main/simple-http/pkg/common/db"
	"simple-main/simple-http/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/01/28 22:17
 @ProductName: video.go
 @Description: 用户视频表数据模型
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
	updateRes := tx.Model(&User{}).Where("id = ?", v.AuthorId).
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

func QueryVideos(ctx context.Context, options ...PageOption) ([]Video, error) {
	res := make([]Video, 0)

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

func QueryVideosByUserId(ctx context.Context, userId int64) ([]Video, error) {
	res := make([]Video, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("author_id = ?", userId).Order("created_at DESC").Find(&res).Error; err != nil {
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

func UpdateVideoCoverUri(ctx context.Context, video *Video) error {
	return db.GetInstance().WithContext(ctx).Model(video).Update("cover_uri", video.CoverUri).Error
}
