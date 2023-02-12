package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple-main/simple-http/cmd/common/db"
	"simple-main/simple-http/cmd/configs"
)

/*
 @Author: 71made
 @Date: 2023/02/01 15:49
 @ProductName: comment.go
 @Description:
*/

type Comment struct {
	gorm.Model
	UserId  uint
	VideoId uint
	Content string
}

func (c *Comment) TableName() string {
	return configs.CommentTable
}

// BeforeCreate
// 通过 GORM 提供的 Hook 实现关联更新 video 记录的 comment_count
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	return c.syncUpdateCommentCount(tx, gorm.Expr("comment_count + 1"))
}

// BeforeDelete
// 同理, 通过 Hook 实现关联更新 video 记录的 comment_count
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	return c.syncUpdateCommentCount(tx, gorm.Expr("comment_count - 1"))
}

func (c *Comment) syncUpdateCommentCount(tx *gorm.DB, expr clause.Expr) (err error) {
	updateRes := tx.Model(&Video{}).Where("id = ?", c.VideoId).
		Update("comment_count", expr)
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

func CreateComment(ctx context.Context, c *Comment) error {
	return db.GetInstance().WithContext(ctx).Create(c).Error
}

func DeleteComment(ctx context.Context, commentId int64) error {
	deleteRes := db.GetInstance().WithContext(ctx).Delete(&Comment{}, commentId)
	if deleteRes.Error != nil {
		return deleteRes.Error
	}
	if deleteRes.RowsAffected <= 0 {
		return errors.New("delete video_comment record fail")
	}
	if deleteRes.RowsAffected > 1 {
		// 同样对删除影响结果做兜底处理
		return errors.New("video_comment table records is dirty")
	}
	return nil
}

func QueryComments(ctx context.Context, videoId int64) ([]Comment, error) {
	res := make([]Comment, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("video_id = ?", videoId).
		Order("created_at DESC").
		Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
