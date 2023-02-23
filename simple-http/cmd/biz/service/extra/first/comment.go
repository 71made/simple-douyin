package first

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/model"
	"simple-main/simple-http/pkg/configs"
	"sync"
)

/*
 @Author: 71made
 @Date: 2023/01/31 15:42
 @ProductName: comment.go
 @Description:
*/

const (
	PublishComment = 1
	RemoveComment  = 2
)

type CommentRequest struct {
	UserId     int64  `query:"user_id"`
	VideoId    int64  `query:"video_id,required"`
	ActionType int    `query:"action_type,required"`
	Content    string `query:"comment_text"`
	CommentId  int64  `query:"comment_id"`
}

type CommentResponse struct {
	biz.Response
	Comment *biz.Comment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	biz.Response
	CommentList []biz.Comment `json:"comment_list,omitempty"`
}

type CommentService interface {
	Action(ctx context.Context, req *CommentRequest) (resp *CommentResponse)
	CommentList(ctx context.Context, userId, videoId int64) (resp *CommentListResponse)
}

func CommentServiceImpl() CommentService {
	return csInstance
}

// commentServiceImpl 对应服务接口实现
type commentServiceImpl struct{}

var csInstance = &commentServiceImpl{}

func (cs *commentServiceImpl) Action(ctx context.Context, req *CommentRequest) (resp *CommentResponse) {
	resp = &CommentResponse{}

	switch req.ActionType {
	case PublishComment:
		{
			newComment, err := cs.publishComment(ctx, req)
			if err != nil {
				hlog.Error(err)
				resp.Response = *biz.NewErrorResponse(err)
				return
			}
			resp.Response = *biz.NewSuccessResponse("评论成功")
			resp.Comment = newComment
		}
	case RemoveComment:
		{
			if err := model.DeleteComment(ctx, req.CommentId, req.VideoId); err != nil {
				hlog.Error(err)
				resp.Response = *biz.NewErrorResponse(err)
				return
			}
			resp.Response = *biz.NewSuccessResponse("删除成功")
		}
	default:
		resp.Response = *biz.NewFailureResponse("非法操作")
	}
	return
}

func (cs *commentServiceImpl) CommentList(ctx context.Context, thisUserId, videoId int64) (resp *CommentListResponse) {
	resp = &CommentListResponse{}

	comments, err := model.QueryComments(ctx, videoId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		resp.CommentList = make([]biz.Comment, 0)
		return
	}

	// 用于缓存 user 映射关系
	var userMap = make(map[uint]*biz.User)
	var isFollowMap = make(map[uint]bool)

	// 初始化 map
	for _, comment := range comments {
		if _, found := userMap[comment.UserId]; !found {
			userMap[comment.UserId] = nil
			isFollowMap[comment.UserId] = false
		}
	}
	// 保存 userId
	var userIds = make([]int64, 0, len(userMap))
	for id := range userMap {
		userIds = append(userIds, int64(id))
	}

	// 并发查询 user 和关注关系
	var wg sync.WaitGroup
	wg.Add(2)

	var QUserErr, QRelationErr error

	go func() {
		defer wg.Done()

		users, err := model.QueryUsersByIds(ctx, userIds)
		if QUserErr = err; QUserErr != nil {
			return
		}
		for _, user := range users {
			userMap[user.ID] = &biz.User{
				Id:                 int64(user.ID),
				Name:               user.Username,
				AvatarURL:          configs.ServerAddr + configs.AvatarURIPrefix + user.Avatar,
				WorkCount:          user.VideoCount,
				BackgroundImage:    user.BackgroundImage,
				Signature:          user.Signature,
				FavoriteCount:      user.FavoriteCount,
				TotalFavoriteCount: user.TotalFavoriteCount,
				FollowCount:        user.FollowCount,
				FollowerCount:      user.FollowerCount,
			}
		}
	}()

	go func() {
		defer wg.Done()
		// 未登陆用户关注关系默认为 false, 直接返回即可
		if thisUserId == biz.NotLoginUserId {
			return
		}
		relations, err := model.QueryRelations(ctx, thisUserId, userIds)
		if QRelationErr = err; QRelationErr != nil {
			return
		}
		for _, relation := range relations {
			isFollowMap[relation.UserId] = relation.IsFollowing()
		}
	}()

	wg.Wait()

	// 处理异常错误
	if QUserErr != nil {
		hlog.Error(QUserErr)
		resp.Response = *biz.NewErrorResponse(QUserErr)
		return
	}
	if QRelationErr != nil {
		hlog.Error(QRelationErr)
		resp.Response = *biz.NewErrorResponse(QRelationErr)
		return
	}

	// 最终转换
	commentList := make([]biz.Comment, len(comments))
	for i, comment := range comments {
		user := userMap[comment.UserId]
		user.IsFollow = isFollowMap[comment.UserId]
		commentList[i] = biz.Comment{
			Id:         int64(comment.ID),
			User:       *user,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		}
	}

	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.CommentList = commentList
	return
}

func (cs *commentServiceImpl) publishComment(ctx context.Context, req *CommentRequest) (*biz.Comment, error) {
	newComment := &model.Comment{
		UserId:  uint(req.UserId),
		VideoId: uint(req.VideoId),
		Content: req.Content,
	}

	// 并发控制
	var wg sync.WaitGroup
	wg.Add(2)

	// 并发创建和查询
	var createErr error
	go func() {
		defer wg.Done()
		if createErr = model.CreateComment(ctx, newComment); createErr != nil {
			return
		}
	}()

	var user *biz.User
	var queryErr error
	go func() {
		defer wg.Done()
		var u *model.User
		u, queryErr = model.QueryUserById(ctx, req.UserId)
		if queryErr != nil {
			return
		}
		user = &biz.User{
			Id:                 int64(u.ID),
			Name:               u.Username,
			AvatarURL:          configs.ServerAddr + configs.AvatarURIPrefix + u.Avatar,
			WorkCount:          u.VideoCount,
			BackgroundImage:    u.BackgroundImage,
			Signature:          u.Signature,
			FavoriteCount:      u.FavoriteCount,
			TotalFavoriteCount: u.TotalFavoriteCount,
			FollowCount:        u.FollowCount,
			FollowerCount:      u.FollowerCount,
			IsFollow:           false, // 对于用户自己, IsFollow 实际上就是默认的 false
		}
	}()

	wg.Wait()

	// 处理异常错误
	if createErr != nil {
		return nil, createErr
	}
	if queryErr != nil {
		return nil, queryErr
	}

	return &biz.Comment{
		Id:         int64(newComment.ID),
		User:       *user,
		Content:    newComment.Content,
		CreateDate: newComment.CreatedAt.Format("01-02"),
	}, nil
}
