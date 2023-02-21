package dal

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple-main/http-rcp/cmd/favorite/dal/db"
	"simple-main/http-rcp/pkg/configs"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/02/21 01:39
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
	FavoriteType uint `gorm:"column:is_favorite"` // 1-点赞 2-取消点赞
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
// 通过 GORM 提供的 Hook 实现关联更新 video 和 user 记录的 favorite_count
func (f *Favorite) BeforeCreate(tx *gorm.DB) (err error) {
	err = f.syncUpdateFavoriteCount(tx, gorm.Expr("`favorite_count` + 1"))
	if err != nil {
		return
	}
	return f.syncUpdateAuthorTotalFavoriteCount(tx, gorm.Expr("`total_favorite_count` + 1"))
}

// BeforeUpdate
// 同理, 通过 Hook 实现关联更新 video 和 user 记录的 favorite_count
func (f *Favorite) BeforeUpdate(tx *gorm.DB) (err error) {
	var fcExpr, tfcExpr clause.Expr
	if f.IsFavorite() {
		fcExpr = gorm.Expr("`favorite_count` + 1")
		tfcExpr = gorm.Expr("`total_favorite_count` + 1")
	} else {
		fcExpr = gorm.Expr("`favorite_count` - 1")
		tfcExpr = gorm.Expr("`total_favorite_count` - 1")
	}

	err = f.syncUpdateFavoriteCount(tx, fcExpr)
	if err != nil {
		return
	}
	return f.syncUpdateAuthorTotalFavoriteCount(tx, tfcExpr)
}

func (f *Favorite) syncUpdateFavoriteCount(tx *gorm.DB, expr clause.Expr) (err error) {
	err = f.syncUpdateVideoFavoriteCount(tx, expr)
	if err != nil {
		return
	}
	return f.syncUpdateUserFavoriteCount(tx, expr)
}

func (f *Favorite) syncUpdateVideoFavoriteCount(tx *gorm.DB, expr clause.Expr) (err error) {

	updateRes := tx.Table(configs.VideoTable).Where("id = ?", f.VideoId).
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

func (f *Favorite) syncUpdateUserFavoriteCount(tx *gorm.DB, expr clause.Expr) (err error) {
	updateRes := tx.Table(configs.UserTable).Where("id = ?", f.UserId).
		Update("favorite_count", expr)
	if err = updateRes.Error; err != nil {
		return err
	}
	if updateRes.RowsAffected <= 0 {
		return errors.New("update user record fail")
	}
	if updateRes.RowsAffected > 1 {
		// 做兜底处理
		return errors.New("user table records is dirty")
	}
	return nil
}

func (f *Favorite) syncUpdateAuthorTotalFavoriteCount(tx *gorm.DB, expr clause.Expr) (err error) {
	var authorId uint
	if err = tx.Table(configs.VideoTable).
		Select("author_id").
		Where("id = ?", f.VideoId).
		Take(&authorId).Error; err != nil {
		// 此处当没找到记录需要抛出错误回滚事务
		return err
	}

	updateRes := tx.Table(configs.UserTable).Where("id = ?", authorId).
		Update("total_favorite_count", expr)
	if err = updateRes.Error; err != nil {
		return err
	}
	if updateRes.RowsAffected <= 0 {
		return errors.New("update user record fail")
	}
	if updateRes.RowsAffected > 1 {
		// 做兜底处理
		return errors.New("user table records is dirty")
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

func QueryFavorites(ctx context.Context, userId int64, videoIds []int64) ([]*Favorite, error) {
	res := make([]*Favorite, 0)

	if len(videoIds) == 0 {
		return res, nil
	}

	if err := db.GetInstance().WithContext(ctx).
		Model(&Favorite{}).
		Where("user_id = ?", userId).
		Where("video_id in ?", videoIds).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryUserFavorites(ctx context.Context, userId int64) ([]*Favorite, error) {
	res := make([]*Favorite, 0)
	if err := db.GetInstance().WithContext(ctx).
		Model(&Favorite{}).
		Where("user_id = ?", userId).
		Where("is_favorite = ?", Favorable).
		Order("updated_at DESC").
		Find(&res).Error; err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return make([]*Favorite, 0), nil
	}

	return res, nil
}
