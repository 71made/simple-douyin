package pack

import (
	"simple-main/http-rcp/cmd/comment/dal"
	csvr "simple-main/http-rcp/grpc_gen/comment"
)

/*
 @Author: 71made
 @Date: 2023/02/21 17:38
 @ProductName: comment.go
 @Description:
*/

func Comment(c *dal.Comment) *csvr.Comment {
	if c == nil {
		return nil
	}

	return &csvr.Comment{
		Id:        int64(c.ID),
		UserId:    int64(c.UserId),
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Unix(),
	}
}

func Comments(cs []*dal.Comment) []*csvr.Comment {
	if cs == nil || len(cs) == 0 {
		return make([]*csvr.Comment, 0)
	}

	res := make([]*csvr.Comment, len(cs))
	for i, c := range cs {
		res[i] = Comment(c)
	}

	return res
}
