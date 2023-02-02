package first

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/cmd/biz"
	"simple-main/cmd/model"
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
	UserId     int64  `json:"user_id"`
	VideoId    int64  `json:"video_id"`
	ActionType int    `json:"action_type"`
	Content    string `json:"comment_content,omitempty"`
	CommentId  int64  `json:"comment_id,omitempty"`
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
		cs.publishComment(ctx, req, resp)
	case RemoveComment:
		cs.removeComment(ctx, req, resp)
	default:
		resp.Response = *biz.NewFailureResponse("未知操作")
	}
	return
}

func (cs *commentServiceImpl) CommentList(ctx context.Context, userId, videoId int64) (resp *CommentListResponse) {
	resp = &CommentListResponse{}

	comments, err := model.QueryComments(ctx, videoId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		resp.CommentList = make([]biz.Comment, 0)
		return
	}

	commentList := make([]biz.Comment, len(comments))

	// 用于缓存 user 映射关系
	var userMap = make(map[uint]*model.User)

	for _, comment := range comments {
		if _, found := userMap[comment.UserId]; !found {
			userMap[comment.UserId] = nil
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

	var queryErr error

	go func() {
		users, err := model.QueryUsersByIds(ctx, userIds)
		if queryErr != nil {
			queryErr = err
			wg.Done()
			return
		}
		for _, user := range users {
			userMap[user.ID] = &user
		}
		wg.Done()
	}()

	go func() {
		wg.Done()
	}()

	wg.Wait()
	// 最终转换
	for i, comment := range comments {
		user := userMap[comment.UserId]
		commentList[i] = biz.Comment{
			Id: int64(comment.ID),
			User: biz.User{
				Id:            int64(user.ID),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				//IsFollow:      false,
			},
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		}
	}

	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.CommentList = commentList
	return
}

func (cs *commentServiceImpl) publishComment(ctx context.Context, req *CommentRequest, resp *CommentResponse) {
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
		if createErr = model.CreateComment(ctx, newComment); createErr != nil {
			hlog.Error(createErr)
		}
		wg.Done()
	}()

	var user biz.User
	var queryErr error
	go func() {
		var u *model.User
		u, queryErr = model.QueryUserById(ctx, req.UserId)
		if queryErr != nil {
			hlog.Error(queryErr)
		}
		user = biz.User{
			Id:            int64(u.ID),
			Name:          u.Username,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			//IsFollow:      false,
		}
		wg.Done()
	}()

	wg.Wait()
	if createErr != nil {
		resp.Response = *biz.NewErrorResponse(createErr)
		return
	}
	if queryErr != nil {
		resp.Response = *biz.NewErrorResponse(queryErr)
		return
	}

	resp.Response = *biz.NewSuccessResponse("评论成功")
	resp.Comment = &biz.Comment{
		Id:         int64(newComment.ID),
		User:       user,
		Content:    newComment.Content,
		CreateDate: newComment.CreatedAt.Format("01-02"),
	}
}
func (cs *commentServiceImpl) removeComment(ctx context.Context, req *CommentRequest, resp *CommentResponse) {
	if err := model.DeleteComment(ctx, req.CommentId); err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	resp.Response = *biz.NewSuccessResponse("删除成功")
}
