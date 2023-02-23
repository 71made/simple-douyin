package pack

import (
	"simple-main/http-rcp/cmd/relation/dal"
	rsvr "simple-main/http-rcp/grpc_gen/relation"
)

/*
 @Author: 71made
 @Date: 2023/02/22 23:00
 @ProductName: relation.go
 @Description:
*/

func Relation(r *dal.Relation) *rsvr.Relation {
	if r != nil {
		return nil
	}

	return &rsvr.Relation{
		Id:          int64(r.ID),
		UserId:      int64(r.UserId),
		FollowerId:  int64(r.FollowerId),
		IsFollowing: r.IsFollowing(),
		IsFriend:    r.IsFriend(),
	}
}

func Relations(rs []*dal.Relation) []*rsvr.Relation {
	if rs == nil || len(rs) == 0 {
		return make([]*rsvr.Relation, 0)
	}

	res := make([]*rsvr.Relation, len(rs))
	for i, r := range rs {
		res[i] = Relation(r)
	}

	return res
}
