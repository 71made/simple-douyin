package second

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"simple-main/simple-http/cmd/biz"
	"simple-main/simple-http/cmd/model"
)

/*
 @Author: 71made
 @Date: 2023/02/03 17:07
 @ProductName: message.go
 @Description:
*/

const (
	Send = 1
)

type MessageRequest struct {
	UserId     int64  `query:"user_id"`
	ToUserId   int64  `query:"to_user_id,required"`
	ActionType int    `query:"action_type,required"`
	Content    string `query:"content,required"`
}

type MessageChatResponse struct {
	biz.Response
	MessageList []biz.Message `json:"message_list"`
}

type MessageService interface {
	Action(ctx context.Context, req *MessageRequest) (resp *biz.Response)
	Chat(ctx context.Context, fromUserId, toUserId int64) (resp *MessageChatResponse)
}

func MessageServiceImpl() MessageService {
	return msInstance
}

type messageServiceImpl struct{}

var msInstance = &messageServiceImpl{}

func (ms *messageServiceImpl) Action(ctx context.Context, req *MessageRequest) (resp *biz.Response) {
	switch req.ActionType {
	case Send:
		{
			if req.UserId == req.ToUserId {
				resp = biz.NewFailureResponse("不可以给自己发送消息")
				return
			}
			// 构建实体
			newMessage := &model.Message{
				FromUserId: uint(req.UserId),
				ToUserId:   uint(req.ToUserId),
				Content:    req.Content,
			}

			if err := model.CreateMessage(ctx, newMessage); err != nil {
				hlog.Error(err)
				resp = biz.NewErrorResponse(err)
				return
			}

			resp = biz.NewSuccessResponse("发送消息成功")
		}
	default:
		resp = biz.NewFailureResponse("非法操作")
	}
	return
}

func (ms *messageServiceImpl) Chat(ctx context.Context, fromUserId, toUserId int64) (resp *MessageChatResponse) {
	resp = &MessageChatResponse{}

	// 对于未登陆用户, 返回空列表即可
	if fromUserId == biz.NotLoginUserId {
		resp.Response = *biz.NewSuccessResponse("获取成功")
		resp.MessageList = make([]biz.Message, 0)
		return
	}

	messages, err := model.QueryMessages(ctx, fromUserId, toUserId)
	if err != nil {
		hlog.Error(err)
		resp.Response = *biz.NewErrorResponse(err)
		return
	}

	// 转换
	messageList := make([]biz.Message, len(messages))
	for i, message := range messages {
		messageList[i] = biz.Message{
			Id:         int64(message.ID),
			FromUserId: int64(message.FromUserId),
			ToUserId:   int64(message.ToUserId),
			Content:    message.Content,
			CreateTime: message.CreatedAt.Unix(),
		}
	}
	resp.Response = *biz.NewSuccessResponse("获取成功")
	resp.MessageList = messageList
	return
}
