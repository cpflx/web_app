package mysql

import (
	"database/sql"

	"go.uber.org/zap"

	"web_app/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("communityList is empty")
			err = nil
		}
	}
	return

}

func GetCommunityDetail(id int64) (community *models.Community, err error) {
	sqlStr := `select community_id, community_name from community where community_id = ?`
	community = &models.Community{}
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("community is null")
			err = nil
		}
	}
	return

}
