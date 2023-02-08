package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"simple-main/cmd/common/db"
	"simple-main/cmd/configs"
	"strings"
)

/*
 @Author: 71made
 @Date: 2023/02/03 17:38
 @ProductName: message.go
 @Description:
*/

const (
	// SendMsgType 当前请求用户发送的消息
	SendMsgType = 1
	// ReceiveMsgType 当前请求用户接收的消息
	ReceiveMsgType = 0

	UnknownMsgType = -1
)

type Message struct {
	gorm.Model
	FromUserId uint
	ToUserId   uint
	Content    string
}

func (m *Message) TableName() string {
	return configs.MessageTable
}

func (m *Message) GetMsgType(userId uint) int {
	if m.FromUserId == userId {
		return SendMsgType
	}
	if m.ToUserId == userId {
		return ReceiveMsgType
	}
	return UnknownMsgType
}

func CreateMessage(ctx context.Context, m *Message) error {
	return db.GetInstance().WithContext(ctx).Create(m).Error
}

func DeleteMessage(ctx context.Context, messageId, userId int64) error {
	deleteRes := db.GetInstance().WithContext(ctx).
		Where("from_user_id = ?", userId).
		Or("to_user_id = ?", userId).
		Delete(&Message{}, messageId)
	if deleteRes.Error != nil {
		return deleteRes.Error
	}

	if deleteRes.RowsAffected <= 0 {
		return errors.New("delete message record fail")
	}
	if deleteRes.RowsAffected > 1 {
		// 做兜底处理
		return errors.New("message table records is dirty")
	}
	return nil
}

func QueryMessages(ctx context.Context, fromId, toId int64) ([]Message, error) {
	res := make([]Message, 0)
	if err := db.GetInstance().WithContext(ctx).
		Where("(from_user_id = ? and to_user_id = ?)", fromId, toId).
		Or("(to_user_id = ? and from_user_id = ?)", fromId, toId).
		Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func QueryFriendMessages(ctx context.Context, userIds []int64) ([]Message, error) {
	res := make([]Message, len(userIds))
	// 构造 SQL
	selects := make([]*gorm.DB, len(userIds))
	var sqlBuilder strings.Builder
	for i, id := range userIds {
		tx := db.GetInstance().WithContext(ctx).Model(&Message{}).Select("*").
			Where("from_user_id = ?", id).
			Or("to_user_id = ?", id).
			Order("created_at DESC").
			Limit(1)
		selects[i] = tx
		sqlBuilder.WriteString("? UNION")
	}
	sql := strings.TrimSuffix(sqlBuilder.String(), " UNION")
	if len(sql) == 0 {
		return res, nil
	}
	// 查询
	if err := db.GetInstance().WithContext(ctx).Raw(sql, selects).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
