package second

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/configs"
	"simple-main/simple-http/cmd/model"
	"sync"
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
	FollowList(ctx context.Context, userId, thisUserId int64) (resp *RelationListResponse)
	FollowerList(ctx context.Context, userId, thisUserId int64) (resp *RelationListResponse)
	FriendList(ctx context.Context, userId int64) (resp *RelationFriendListResponse)
}

func RelationServiceImpl() RelationService {
	return rsInstance
}

type relationServiceImpl struct{}

var rsInstance = &relationServiceImpl{}

func (rs *relationServiceImpl) Action(ctx context.Context, req *RelationActionRequest) (resp *biz.Response) {

	if req.UserId == req.ToUserId {
		resp = biz.NewFailureResponse("不能关注自己")
		return
	}

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

func (rs *relationServiceImpl) FollowList(ctx context.Context, userId, thisUserId int64) (resp *RelationListResponse) {
	resp = &RelationListResponse{}

	relations, err := model.QueryFollowRelations(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	// 是否是 thisUser 的 follow list
	isFollow := userId == thisUserId
	userList, err := rs.transToUsers(ctx, relations, thisUserId, isFollow)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.UserList = userList
	return
}

func (rs *relationServiceImpl) FollowerList(ctx context.Context, userId, thisUserId int64) (resp *RelationListResponse) {
	resp = &RelationListResponse{}

	relations, err := model.QueryFollowerRelations(ctx, userId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	userList, err := rs.transToUsers(ctx, relations, thisUserId, false)
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

func (rs *relationServiceImpl) transToUsers(ctx context.Context, relations []model.Relation, thisUserId int64, isFollow bool) ([]biz.User, error) {
	if len(relations) == 0 {
		return make([]biz.User, 0), nil
	}
	userIds := make([]int64, len(relations))

	// 用于缓存 isFollow 映射关系
	var isFollowMap = make(map[uint]bool, len(relations))
	for i, relation := range relations {
		// 对于 thisUser 自己的关注列表, isFollow 直接为 true
		if isFollow {
			userIds[i] = int64(relation.UserId)
			isFollowMap[relation.UserId] = true
		} else {
			userIds[i] = int64(relation.FollowerId)
			isFollowMap[relation.FollowerId] = false

		}
	}

	// 最后构造 List
	userList := make([]biz.User, len(userIds))

	// 并发控制
	var wg sync.WaitGroup
	wg.Add(2)
	// 并发查询
	var QUsersErr, QRelationsErr error
	var users []model.User

	go func() {
		users, QUsersErr = model.QueryUsersByIds(ctx, userIds)
		wg.Done()
	}()

	go func() {
		defer wg.Done()
		// 对于非 thisUser 关注列表, 需要查询对应关系
		if !isFollow && thisUserId != biz.NotLoginUserId {
			var reverseRelations []model.Relation
			reverseRelations, QRelationsErr = model.QueryRelations(ctx, thisUserId, userIds)
			if QRelationsErr != nil {
				return
			}
			for _, relation := range reverseRelations {
				isFollowMap[relation.UserId] = relation.IsFollowing()
			}
		}
	}()

	wg.Wait()
	// 处理异常错误
	if QUsersErr != nil {
		return nil, QUsersErr
	}
	if QRelationsErr != nil {
		return nil, QRelationsErr
	}

	for i, user := range users {
		userList[i] = biz.User{
			Id:            int64(user.ID),
			Name:          user.Username,
			AvatarURL:     configs.ServerAddr + configs.AvatarURIPrefix + user.Avatar,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollowMap[user.ID],
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

	// 并发控制
	var wg sync.WaitGroup
	wg.Add(2)

	// 并发查询
	var QUserErr, QMessageErr error
	var users []model.User
	var messages []model.Message

	go func() {
		users, QUserErr = model.QueryUsersByIds(ctx, userIds)
		wg.Done()
	}()

	go func() {
		// 查询朋友第一条消息
		messages, QMessageErr = model.QueryFriendMessages(ctx, userIds)
		wg.Done()
	}()

	wg.Wait()
	// 处理异常错误
	if QUserErr != nil {
		return nil, QUserErr
	}
	if QMessageErr != nil {
		return nil, QMessageErr
	}

	// 最后拼接
	userList := make([]biz.FriendUser, len(users))
	for i, user := range users {
		userList[i] = biz.FriendUser{
			// 朋友列表中必定是已关注用户, 直接设置 IsFollow 为 true
			// 但不一定是互关用户, 可以存在对方已经取消关注
			User: biz.User{
				Id:            int64(user.ID),
				Name:          user.Username,
				AvatarURL:     configs.ServerAddr + configs.AvatarURIPrefix + user.Avatar,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      true,
			},
		}
		userMap[user.ID] = &userList[i]
	}
	// 设置消息
	for _, message := range messages {
		// message 不为零值, 即查询到记录
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
