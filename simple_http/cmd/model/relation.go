package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple-main/simple-http/cmd/common/db"
	"simple-main/simple-http/cmd/configs"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/02/04 17:44
 @ProductName: relation.go
 @Description:
*/

const (
	Following    = 1
	NotFollowing = 2
)

const (
	IsFriend  = 1
	NotFriend = 2
)

type Relation struct {
	ID           uint `gorm:"primarykey"`
	UserId       uint
	FollowerId   uint
	FollowType   uint `gorm:"column:is_following"` // 1-关注 2-取消关注
	FriendStatus uint `gorm:"column:is_friend"`    // 1-朋友 2-不是朋友
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *Relation) TableName() string {
	return configs.RelationTable
}

func (r *Relation) IsFollowing() bool {
	return r.FollowType == Following
}

func (r *Relation) IsFriend() bool {
	return r.FriendStatus == IsFriend
}

func (r *Relation) GetFollowType() uint {
	// 过滤其他参数
	if r.FollowType != Following {
		r.FollowType = NotFollowing
	}
	return r.FollowType
}

func (r *Relation) GetFriendStatus() uint {
	// 过滤其他参数
	if r.FriendStatus != IsFriend {
		r.FriendStatus = NotFriend
	}
	return r.FriendStatus
}

func (r *Relation) BeforeSave(tx *gorm.DB) (err error) {
	// 处理朋友关系
	return r.updateFriendShip(tx)
}

// BeforeCreate
// 通过 GORM 提供的 Hook 实现关联更新 user 记录的 follow_count 和 follower_count
func (r *Relation) BeforeCreate(tx *gorm.DB) (err error) {

	// 更新关注数和粉丝数
	followerExpr := gorm.Expr("follower_count + 1")
	followExpr := gorm.Expr("follow_count + 1")

	if err = r.syncUpdateFollowerCount(tx, followerExpr); err != nil {
		return err
	}
	if err = r.syncUpdateFollowCount(tx, followExpr); err != nil {
		return err
	}
	return
}

// BeforeUpdate
// 同理, 通过 Hook 实现关联更新 user 记录的 follow_count 和 follower_count
func (r *Relation) BeforeUpdate(tx *gorm.DB) (err error) {

	// 更新关注数和粉丝数
	var followerExpr, followExpr clause.Expr
	if r.IsFollowing() {
		followerExpr = gorm.Expr("follower_count + 1")
		followExpr = gorm.Expr("follow_count + 1")
	} else {
		followerExpr = gorm.Expr("follower_count - 1")
		followExpr = gorm.Expr("follow_count - 1")
	}

	if err = r.syncUpdateFollowerCount(tx, followerExpr); err != nil {
		return err
	}
	if err = r.syncUpdateFollowCount(tx, followExpr); err != nil {
		return err
	}
	return
}

func (r *Relation) updateFriendShip(tx *gorm.DB) (err error) {
	// 对于 FriendStatus 主动设置为 NotFriend, 认为 r 是更新实体, 直接交由 UpdateRelation 方法处理
	if r.FriendStatus == NotFriend {
		return
	}
	if r.IsFollowing() {
		// 对于关注操作, 尝试更新 互关用户的 relation
		updateRes := tx.Table(r.TableName()).
			Where("user_id = ?", r.FollowerId).
			Where("follower_id = ?", r.UserId).
			Where("is_following = ?", Following).
			Update("is_friend", IsFriend)

		if updateRes.Error != nil {
			err = updateRes.Error
			return
		}
		if updateRes.RowsAffected == 0 {
			// 没有匹配记录则设置 NotFriend
			r.FriendStatus = NotFriend
			tx.Statement.SetColumn("is_friend", NotFriend)
		} else if updateRes.RowsAffected == 1 {
			// 有匹配记录或者修改成功则设置为 IsFriend
			// 对于 MySQL 此处需要保证数据库连接 DSN 中配置了 clientFoundRows 为 true
			// 因为 MySQL 连接驱动默认只返回影响行数, 不会返回匹配行数
			// 详细参考: https://my.oschina.net/zdtdtel/blog/5321128
			r.FriendStatus = IsFriend
			tx.Statement.SetColumn("is_friend", IsFriend)
		} else {
			// 做兜底处理
			return errors.New("user_video table records is dirty")
		}
	} else {
		// 取消关注操作, 则修改 FriendStatus 为 NotFriend 即可
		// 实际上, 正常流程在调用取消关注接口时, service 层会设置 FriendStatus 为 NotFriend
		// 并在此函数开始被直接返回, 这里做兜底处理, 处理 FriendStatus 零值时的情况
		r.FriendStatus = NotFriend
		tx.Statement.SetColumn("is_friend", NotFriend)
	}
	return
}

func (r *Relation) syncUpdateFollowerCount(tx *gorm.DB, expr clause.Expr) (err error) {
	updateRes := tx.Model(&User{}).Where("id = ?", r.UserId).
		Update("follower_count", expr)
	return r.syncErrors(updateRes)
}

func (r *Relation) syncUpdateFollowCount(tx *gorm.DB, expr clause.Expr) (err error) {
	updateRes := tx.Model(&User{}).Where("id = ?", r.FollowerId).
		Update("follow_count", expr)
	return r.syncErrors(updateRes)
}

func (r *Relation) syncErrors(updateRes *gorm.DB) error {
	if err := updateRes.Error; err != nil {
		return err
	}
	if updateRes.RowsAffected <= 0 {
		return errors.New("update user_video record fail")
	}
	if updateRes.RowsAffected > 1 {
		// 做兜底处理
		return errors.New("user_video table records is dirty")
	}
	return nil
}

func CreateRelation(ctx context.Context, r *Relation) error {
	return db.GetInstance().WithContext(ctx).Create(r).Error
}

func UpdateRelation(ctx context.Context, r *Relation) error {
	err := db.GetInstance().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updateRes := tx.Model(r).
			Where("user_id", r.UserId).Where("follower_id", r.FollowerId).
			Updates(Relation{
				FollowType:   r.GetFollowType(),
				FriendStatus: r.GetFriendStatus(),
			})

		if updateRes.Error != nil {
			return updateRes.Error
		}
		if updateRes.RowsAffected <= 0 {
			return errors.New("update relation record fail")
		}
		if updateRes.RowsAffected > 1 {
			// 做兜底处理
			return errors.New("relation table records is dirty")
		}
		return nil
	})

	return err
}

// QueryRelation
// 查询 userId 对于 toUserId 的关注关系, 即 userId 是否关注了 toUserId
func QueryRelation(ctx context.Context, userId, toUserId int64) (*Relation, error) {
	res := make([]Relation, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("user_id = ?", toUserId).
		Where("follower_id = ?", userId).
		Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func QueryRelations(ctx context.Context, userId int64, toUserIds []int64) ([]Relation, error) {
	res := make([]Relation, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("user_id in ?", toUserIds).
		Where("follower_id = ?", userId).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func QueryFollowRelations(ctx context.Context, userId int64) ([]Relation, error) {
	res := make([]Relation, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("follower_id = ?", userId).
		Where("is_following = ?", Following).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func QueryFollowerRelations(ctx context.Context, toUserId int64) ([]Relation, error) {
	res := make([]Relation, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("user_id = ?", toUserId).
		Where("is_following = ?", Following).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func QueryFriendRelations(ctx context.Context, userId int64) ([]Relation, error) {
	res := make([]Relation, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("follower_id = ?", userId).
		Where("is_following = ?", Following).
		Where("is_friend = ?", IsFriend).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
