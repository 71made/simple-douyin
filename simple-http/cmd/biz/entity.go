package biz

/*
 @Author: 71made
 @Date: 2023/01/24 22:53
 @ProductName: entity.go
 @Description: 业务实体模型
*/

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

const (
	SuccessCode = iota        // 请求处理成功状态码
	FailureCode               // 请求处理失败状态码
	ErrorCode   = FailureCode // 请求处理错误状态码
)

func NewErrorResponse(err error) *Response {
	return &Response{
		StatusCode: ErrorCode,
		StatusMsg:  err.Error(),
	}
}

func NewSuccessResponse(successMsg string) *Response {
	return &Response{
		StatusCode: SuccessCode,
		StatusMsg:  successMsg,
	}
}

func NewFailureResponse(failureMsg string) *Response {
	return &Response{
		StatusCode: FailureCode,
		StatusMsg:  failureMsg,
	}
}

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
}

type Comment struct {
	Id         int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

// NotLoginUserId 定义未登陆用户 id 为 -1
const NotLoginUserId = -1

type User struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	AvatarURL          string `json:"avatar"`
	BackgroundImage    string `json:"background_image"`
	Signature          string `json:"signature"`
	FollowCount        int64  `json:"follow_count"`
	FollowerCount      int64  `json:"follower_count"`
	WorkCount          int64  `json:"work_count"`
	TotalFavoriteCount int64  `json:"total_favorited"`
	FavoriteCount      int64  `json:"favorite_count"`
	IsFollow           bool   `json:"is_follow"`
}

type Message struct {
	Id         int64  `json:"id"`
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type FriendUser struct {
	User
	Message string `json:"message"`
	MsgType int    `json:"msg_type"`
}
