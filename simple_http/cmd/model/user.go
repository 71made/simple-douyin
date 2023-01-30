package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"simple-main/cmd/common/db"
	"simple-main/pkg/configs"
)

/*
 @Author: 71made
 @Date: 2023/01/24 23:05
 @ProductName: user.go
 @Description: 用户表数据模型
*/

type User struct {
	gorm.Model
	Username      string
	Password      string
	FollowCount   int64
	FollowerCount int64
}

func (u *User) TableName() string {
	return configs.UserTable
}

func QueryUsers(ctx context.Context, username string) ([]User, error) {
	res := make([]User, 0)
	if err := db.GetInstance().WithContext(ctx).Where("username = ?", username).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryUser(ctx context.Context, username string) (*User, error) {
	res, err := QueryUsers(ctx, username)
	if err != nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

func QueryUserById(ctx context.Context, userId int64) (*User, error) {
	res := make([]User, 0)
	if err := db.GetInstance().WithContext(ctx).Where("id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func IsExistUser(ctx context.Context, username string) (bool, error) {
	if len(username) == 0 {
		return false, errors.New("username can't be empty")
	}
	ids := make([]int64, 0)
	if err := db.GetInstance().WithContext(ctx).
		Select("id").
		Where("username = ?", username).
		Find(&ids).Error; err != nil {
		return false, err
	}
	return len(ids) != 0, nil
}

func CreateUser(ctx context.Context, user *User) error {
	return db.GetInstance().WithContext(ctx).Create(user).Error
}
