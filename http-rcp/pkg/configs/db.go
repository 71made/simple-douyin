package configs

/*
 @Author: 71made
 @Date: 2023/01/24 21:58
 @ProductName: db.go
 @Description: 数据库配置
*/

const (
	MySQLDataBaseDSN = "gorm:gorm@tcp(localhost:3308)/douyin?charset=utf8&parseTime=True&loc=Local&clientFoundRows=true"
	UserTable        = "user"
	VideoTable       = "user_video"
	FavoriteTable    = "favorite"
	CommentTable     = "video_comment"
	MessageTable     = "message"
	RelationTable    = "user_relation"
)
