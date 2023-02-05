package second

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/cmd/biz"
	"simple-main/cmd/model"
)

/*
 @Author: 71made
 @Date: 2023/02/03 17:07
 @ProductName: relation.go
 @Description:
*/

type RelationActionRequest struct {
	UserId     int64 `query:"user_id"`
	ToUserId   int64 `query:"to_user_id,required"`
	ActionType uint  `query:"action_type,required"`
}

type RelationListResponse struct {
	biz.Response
	UserList []biz.User `json:"user_list"`
}

type RelationFriendListResponse struct {
	biz.Response
	UserList []biz.FriendUser `json:"user_list"`
}

type RelationService interface {
	Action(ctx context.Context, req *RelationActionRequest) (resp *biz.Response)
	FollowList(ctx context.Context, userId int64) (resp *RelationListResponse)
	FollowerList(ctx context.Context, userId int64) (resp *RelationListResponse)
	FriendList(ctx context.Context, userId int64) (resp *RelationFriendListResponse)
}

func RelationServiceImpl() RelationService {
	return rsInstance
}

type relationServiceImpl struct{}

var rsInstance = &relationServiceImpl{}

func (rs *relationServiceImpl) Action(ctx context.Context, req *RelationActionRequest) (resp *biz.Response) {
	// 构建实体
	r := &model.Relation{
		UserId:     uint(req.ToUserId),
		FollowerId: uint(req.UserId),
		FollowType: req.ActionType,
	}

	// 查找关系记录
	found, err := model.QueryRelation(ctx, req.UserId, req.ToUserId)
	if err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(err)
		return
	}

	// 不存在就创建, 存在则更新
	if found == nil {
		err = model.CreateRelation(ctx, r)
	} else if found.FollowType != r.FollowType {
		// 并且只对于不同的 type, 才触发更新
		if found.IsFriend() {
			// 对于已经是朋友的 relation, 更新为 NotFriend
			r.FriendStatus = model.NotFriend
		}
		err = model.UpdateRelation(ctx, r)
	}

	if err != nil {
		hlog.Error(err)
		resp = biz.NewErrorResponse(err)
		return
	}

	resp = biz.NewSuccessResponse("操作成功")
	return
}

func (rs *relationServiceImpl) FollowList(ctx context.Context, userId int64) (resp *RelationListResponse) {
	resp = &RelationListResponse{}

	relations, err := model.QueryFollowRelations(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	userList, err := rs.transToUsers(ctx, relations)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.UserList = userList
	return
}

func (rs *relationServiceImpl) FollowerList(ctx context.Context, userId int64) (resp *RelationListResponse) {
	resp = &RelationListResponse{}

	relations, err := model.QueryFollowerRelations(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	userList, err := rs.transToUsers(ctx, relations)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.UserList = userList
	return
}

func (rs *relationServiceImpl) FriendList(ctx context.Context, userId int64) (resp *RelationFriendListResponse) {
	resp = &RelationFriendListResponse{}

	relations, err := model.QueryFriendRelations(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	userList, err := rs.transToFriendUsers(ctx, relations, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.UserList = userList
	return
}

func (rs *relationServiceImpl) transToUsers(ctx context.Context, relations []model.Relation) ([]biz.User, error) {
	userIds := make([]int64, len(relations))
	for i, relation := range relations {
		userIds[i] = int64(relation.UserId)
	}

	users, err := model.QueryUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	userList := make([]biz.User, len(users))
	for i, user := range users {
		userList[i] = biz.User{
			Id:            int64(user.ID),
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      relations[i].IsFollowing(),
		}
	}
	return userList, nil
}

func (rs *relationServiceImpl) transToFriendUsers(ctx context.Context, relations []model.Relation, thisUserId int64) ([]biz.FriendUser, error) {
	userIds := make([]int64, len(relations))
	userMap := make(map[uint]*biz.FriendUser, len(relations))
	for i, relation := range relations {
		userIds[i] = int64(relation.UserId)
		userMap[relation.UserId] = nil
	}

	users, err := model.QueryUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	userList := make([]biz.FriendUser, len(users))
	for i, user := range users {
		userList[i] = biz.FriendUser{
			User: biz.User{
				Id:            int64(user.ID),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      relations[i].IsFollowing(),
			},
		}
		userMap[user.ID] = &userList[i]
	}

	// 查询朋友第一条消息
	messages, err := model.QueryFriendMessages(ctx, userIds)
	for _, message := range messages {
		if message.ID != 0 {
			var friendId uint
			if message.ToUserId != uint(thisUserId) {
				friendId = message.ToUserId
			} else {
				friendId = message.FromUserId
			}
			user, ok := userMap[friendId]
			if ok {
				user.Message = message.Content
				user.MsgType = message.GetMsgType(uint(thisUserId))
			}
		}
	}
	return userList, nil
}
