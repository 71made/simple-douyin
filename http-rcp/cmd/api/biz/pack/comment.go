package pack

import (
	"simple-main/http-rcp/cmd/api/biz"
	csvr "simple-main/http-rcp/grpc_gen/comment"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/02/21 21:02
 @ProductName: comment.go
 @Description:
*/

func BizComment(c *csvr.Comment, u *biz.BaseUser) *biz.Comment {
	if c == nil {
		return nil
	}

	return &biz.Comment{
		Id:         c.Id,
		User:       u,
		Content:    c.Content,
		CreateDate: time.Unix(c.CreatedAt, 0).Format("01-02"),
	}
}
