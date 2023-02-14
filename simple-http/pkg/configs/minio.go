package configs

/*
 @Author: 71made
 @Date: 2023/01/29 22:05
 @ProductName: minio.go
 @Description: MinIO 服务相关配置常量
*/

const (
	Endpoint        = "127.0.0.1:9000"
	AccessKeyID     = "minio"
	SecretAccessKey = "12345678"
	UseSSL          = false
)

const (
	BucketName = "douyin"
	Video      = "/video/"
	Cover      = "/cover/"
	Avatar     = "/avatar/"
)

const ServerAddr = "http://192.168.0.107:9000" // 项目本地测试图床 IP 地址

const (
	VideoURIPrefix  = "/" + BucketName + Video // 视频路径 uri 前缀
	CoverURIPrefix  = "/" + BucketName + Cover // 视频封面路径 uri 前缀
	AvatarURIPrefix = "/" + BucketName + Avatar
	EmptyCoverName  = "empty_cover.jpeg"
	EmptyAvatarName = "empty_avatar.jpeg"
)
