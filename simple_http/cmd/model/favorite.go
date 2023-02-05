package model

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple-main/cmd/common/db"
	"simple-main/cmd/configs"
	"strings"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/01/31 16:36
 @ProductName: favorite.go
 @Description:
*/

const (
	Favorable   = 1
	Unfavorable = 2
)

type Favorite struct {
	ID           uint `gorm:"primarykey"`
	UserId       uint
	VideoId      uint
	FavoriteType uint `gorm:"column:is_favorite"` // 1-guan 2-取消点赞
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (f *Favorite) TableName() string {
	return configs.FavoriteTable
}

func (f *Favorite) IsFavorite() bool {
	return f.FavoriteType == Favorable
}

func (f *Favorite) GetFavoriteType() uint {
	// 过滤其他参数
	if f.FavoriteType != Favorable {
		f.FavoriteType = Unfavorable
	}
	return f.FavoriteType
}

// BeforeCreate
// 通过 GORM 提供的 Hook 实现关联更新 video 记录的 favorite_count
func (f *Favorite) BeforeCreate(tx *gorm.DB) (err error) {
	return f.syncUpdateFavoriteCount(tx, gorm.Expr("favorite_count + 1"))
}

// BeforeUpdate
// 同理, 通过 Hook 实现关联更新 video 记录的 favorite_count
func (f *Favorite) BeforeUpdate(tx *gorm.DB) (err error) {
	var expr clause.Expr
	if f.IsFavorite() {
		expr = gorm.Expr("favorite_count + 1")
	} else {
		expr = gorm.Expr("favorite_count - 1")
	}
	return f.syncUpdateFavoriteCount(tx, expr)
}

func (f *Favorite) syncUpdateFavoriteCount(tx *gorm.DB, expr clause.Expr) (err error) {

	updateRes := tx.Model(&Video{}).Where("id = ?", f.VideoId).
		Update("favorite_count", expr)
	if err = updateRes.Error; err != nil {
		return err
	}
	if updateRes.RowsAffected <= 0 {
		return errors.New("update user_video record fail")
	}
	if updateRes.RowsAffected > 1 {
		// 对于影响数超过 1 的更新, 逻辑上是不合理的, 可能是 video 产生脏数据
		// 实际上, 在主键约束下不可能出现此情况, 仅做兜底处理
		return errors.New("user_video table records is dirty")
	}
	return nil
}

func CreateFavorite(ctx context.Context, f *Favorite) error {

	return db.GetInstance().WithContext(ctx).Create(f).Error
}

func UpdateFavorite(ctx context.Context, f *Favorite) error {
	err := db.GetInstance().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updateRes := tx.Model(f).
			Where("user_id", f.UserId).Where("video_id", f.VideoId).
			Update("is_favorite", f.GetFavoriteType())

		if updateRes.Error != nil {
			return updateRes.Error
		}
		if updateRes.RowsAffected <= 0 {
			return errors.New("update favorite record fail")
		}
		if updateRes.RowsAffected > 1 {
			// 做兜底处理
			return errors.New("favorite table records is dirty")
		}
		return nil
	})

	return err
}

func QueryFavorite(ctx context.Context, userId, videoId int64) (*Favorite, error) {
	res := make([]Favorite, 0)
	if err := db.GetInstance().WithContext(ctx).
		Model(&Favorite{}).
		Where("user_id = ?", userId).
		Where("video_id = ?", videoId).
		Find(&res).Error; err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func QueryFavorites(ctx context.Context, userId int64, videoIds []int64) ([]Favorite, error) {
	res := make([]Favorite, 0)
	if err := db.GetInstance().WithContext(ctx).
		Model(&Favorite{}).
		Where("user_id = ?", userId).
		Where("video_id in ?", videoIds).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryFavoriteVideos(ctx context.Context, userId int64) ([]Video, error) {
	videoIds := make([]int64, 0)
	if err := db.GetInstance().WithContext(ctx).
		Select("video_id").Model(&Favorite{}).
		Where("user_id = ?", userId).
		Where("is_favorite = ?", Favorable).
		Order("updated_at DESC").
		Find(&videoIds).Error; err != nil {
		return nil, err
	}

	if len(videoIds) == 0 {
		return make([]Video, 0), nil
	}

	res := make([]Video, len(videoIds))
	// 构造排序条件
	str := strings.ReplaceAll(fmt.Sprintf("%v", videoIds), " ", ",")
	// 截取中间 id 序列
	str = str[1 : len(str)-1]
	if err := db.GetInstance().WithContext(ctx).
		Model(&Video{}).Where("id in ?", videoIds).
		Order("Field(id," + str + ")").
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
