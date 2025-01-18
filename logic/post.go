package logic

import (
	"go.uber.org/zap"

	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post_id
	p.ID = snowflake.GetID()

	// 2. 将数据保存到mysql中
	return mysql.CreatePost(p)
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	postRes, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed", zap.Error(err))
		return nil, err
	}

	data = new(models.ApiPostDetail)
	// 根据作者ID查询作者信息
	user, err := mysql.GetUserByID(postRes.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
		return nil, err
	}

	// 根据社区ID查询社区详情
	community, err := mysql.GetCommunityDetail(postRes.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
		return nil, err
	}

	data = &models.ApiPostDetail{
		AuthorName: user.Username,
		Post:       postRes,
		Community:  community,
	}
	return
}

func GetPostList(offset, limit int64) (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(offset, limit)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, 2)

	for _, detail := range *postList {
		// 根据作者ID查询作者信息
		user, _ := mysql.GetUserByID(detail.AuthorID)
		// 根据社区ID查询社区详情
		community, _ := mysql.GetCommunityDetail(detail.CommunityID)
		data = append(data, &models.ApiPostDetail{
			AuthorName: user.Username,
			Post:       &detail,
			Community:  community,
		})

	}

	return
}
