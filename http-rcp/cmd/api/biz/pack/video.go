package pack

import (
	"context"
	"errors"
	"simple-main/http-rcp/cmd/api/biz"
	"simple-main/http-rcp/cmd/api/rpc"
	vsvr "simple-main/http-rcp/grpc_gen/video"
	"sync"
)

/*
 @Author: 71made
 @Date: 2023/02/19 02:53
 @ProductName: video.go
 @Description:
*/

// BizVideos
// 类型转换 []*vsvr.Video -> []*biz.Video, 用户未登陆时, userId 传递 NotLoginUserId
func BizVideos(ctx context.Context, videos []*vsvr.Video, userId int64) ([]*biz.Video, error) {

	// 缓存 author 映射关系
	var authors = make(map[int64]*biz.User)
	// 缓存 author id
	var authorIds = make([]int64, len(authors))

	// 缓存 video id 为后续查询关系做准备
	var videoIds = make([]int64, len(videos))
	// 保存 video id 与实体间的映射关系
	var videoMap = make(map[int64]*biz.Video, len(videos))

	// 保存关注情况
	var isFollowMap = make(map[int64]bool)

	// 构建转换实体
	var videoList = make([]*biz.Video, len(videos))
	for i, video := range videos {
		// 保存 authorIds 和初始化 authors、isFollowMap
		_, found := authors[video.AuthorId]
		if !found {
			authors[video.AuthorId] = nil
			authorIds = append(authorIds, video.AuthorId)
			isFollowMap[video.AuthorId] = false
		}

		// 构造 video
		videoList[i] = &biz.Video{
			Id:            video.Id,
			PlayURL:       video.PlayUrl,
			CoverURL:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
		}
		videoIds[i] = video.Id
		videoMap[video.Id] = videoList[i]
	}

	// 并发控制
	var wg sync.WaitGroup
	wg.Add(1) // 后续查询 authors 实体的并发任务

	wg.Add(1)
	//// 对于已登陆用户, 需要查询用户对应的关注和点赞情况
	//if userId != biz.NotLoginUserId {
	//	wg.Add(2) // 接下来查询关注和点赞的并发任务
	//
	//	go func() {
	//		// 查询是否关注
	//		relations, _ := model.QueryRelations(ctx, userId, authorIds)
	//
	//		for _, relation := range relations {
	//			isFollowMap[relation.UserId] = relation.IsFollowing()
	//		}
	//		wg.Done()
	//	}()
	//
	go func() {
		// 查询是否点赞
		favorites, _ := rpc.QueryFavorites(ctx, userId, videoIds)

		for _, favorite := range favorites {
			video := videoMap[favorite.VideoId]
			video.IsFavorite = favorite.IsFavorite
		}
		wg.Done()
	}()
	//
	//}
	var QUserErr error
	go func() {
		defer wg.Done()
		// 批量查询
		users, baseResp := rpc.QueryUsers(ctx, authorIds)
		if baseResp != nil && baseResp.StatusCode != biz.SuccessCode {
			QUserErr = errors.New(baseResp.StatusMsg)
		}
		for _, user := range users {
			author := BizUser(user, false) // 先默认为false
			authors[user.Id] = author
		}
	}()

	wg.Wait()

	// 处理异常错误
	if QUserErr != nil {
		return nil, QUserErr
	}

	// 最后组装
	for i, video := range videos {
		author := authors[video.AuthorId]
		author.IsFollow = isFollowMap[video.AuthorId]
		videoList[i].Author = author
	}

	return videoList, nil
}
