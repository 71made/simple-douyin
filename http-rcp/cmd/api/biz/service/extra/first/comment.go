package first

import (
	"context"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/biz/pack"
	"simple-main/http-rcp/cmd/api/rpc"
	csvr "simple-main/http-rcp/grpc_gen/comment"
	usvr "simple-main/http-rcp/grpc_gen/user"
	"sync"
)

/*
 @Author: 71made
 @Date: 2023/02/21 19:09
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
	*biz.Response
	Comment *biz.Comment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	*biz.Response
	CommentList []*biz.Comment `json:"comment_list,omitempty"`
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
			newComment, baseResp := cs.publishComment(ctx, req)
			resp.Response = baseResp
			resp.Comment = newComment
		}
	case RemoveComment:
		{
			if baseResp := rpc.DeleteComment(ctx, req.CommentId, req.VideoId, req.UserId); baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
				resp.Response = baseResp
				return
			}

			resp.Response = biz.NewSuccessResponse("删除成功")
		}
	default:
		resp.Response = biz.NewFailureResponse("非法操作")
	}
	return
}

func (cs *commentServiceImpl) CommentList(ctx context.Context, thisUserId, videoId int64) (resp *CommentListResponse) {
	resp = &CommentListResponse{}

	comments, baseResp := rpc.QueryComments(ctx, videoId)
	if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
		resp.Response = baseResp
		resp.CommentList = make([]*biz.Comment, 0)
		return
	}

	// 用于缓存 user 映射关系
	var userMap = make(map[int64]*biz.BaseUser)
	//var isFollowMap = make(map[int64]bool)

	// 初始化 map
	for _, comment := range comments {
		if _, found := userMap[comment.UserId]; !found {
			userMap[comment.UserId] = nil
			//isFollowMap[comment.UserId] = false
		}
	}
	// 保存 userId
	var userIds = make([]int64, 0, len(userMap))
	for id := range userMap {
		userIds = append(userIds, id)
	}

	// 并发查询 user 和关注关系
	var wg sync.WaitGroup
	wg.Add(1)

	var QUserRes *biz.Response
	go func() {
		defer wg.Done()

		users, baseResp := rpc.QueryUsers(ctx, userIds)
		if QUserRes = baseResp; QUserRes != nil && QUserRes.StatusCode != biz.SuccessCode {
			return
		}
		for _, user := range users {
			userMap[user.Id] = pack.BizBaseUser(user, false) // 先默认给 false
		}
	}()

	wg.Wait()

	// 处理响应
	if QUserRes != nil && QUserRes.StatusCode != biz.SuccessCode {
		resp.Response = QUserRes
		resp.CommentList = make([]*biz.Comment, 0)
		return
	}

	// 最终转换
	commentList := make([]*biz.Comment, len(comments))
	for i, comment := range comments {
		user := userMap[comment.UserId]
		//user.IsFollow = isFollowMap[comment.UserId]
		commentList[i] = pack.BizComment(comment, user)
	}

	resp.Response = biz.NewSuccessResponse("获取成功")
	resp.CommentList = commentList
	return

}

func (cs *commentServiceImpl) publishComment(ctx context.Context, req *CommentRequest) (*biz.Comment, *biz.Response) {
	if len(req.Content) == 0 {
		return nil, biz.NewFailureResponse("评论内容为空")
	}

	// 并发调用 server 服务查询
	var wg sync.WaitGroup
	wg.Add(2)

	var (
		newComment *csvr.Comment
		user       *usvr.User
	)

	var CCommentRes, QUserRes *biz.Response

	go func() {
		defer wg.Done()
		newComment, CCommentRes = rpc.CreateComment(ctx, req.VideoId, req.UserId, req.Content)

	}()

	go func() {
		defer wg.Done()
		user, QUserRes = rpc.QueryUser(ctx, req.UserId)
	}()

	wg.Wait()

	// 处理响应
	if CCommentRes != nil && CCommentRes.StatusCode != biz.SuccessCode {
		return nil, CCommentRes
	}
	if QUserRes != nil && QUserRes.StatusCode != biz.SuccessCode {
		return nil, QUserRes
	}

	return pack.BizComment(newComment, pack.BizBaseUser(user, false)), // 对于用户自己 isFollow 默认为 false
		biz.NewSuccessResponse("评论提交成功")

}
