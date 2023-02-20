package configs

/*
 @Author: 71made
 @Date: 2023/02/13 18:30
 @ProductName: system.go
 @Description:
*/

const TCP = "tcp"
const ServerIP = "127.0.0.1"

// 各微服务端口配置
const (
	UserServerPort  = 8081
	VideoServerPort = 8082
)

const (
	UserServer  = "user"
	VideoServer = "video"
)

const (
	VideoPathPrefix = "./static/video/" // 上传视频文件相对路径前缀
	CoverPathPrefix = "./static/cover/" // 上传视频封面相对路径前缀
)

const (
	JWTSecretKey = "douyin::JWT"
	IdentityKey  = "JWT::UserId"
)
