package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
)

// VoteForPost 为帖子投票
func VoteForPost(c *gin.Context, userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Any("userID", userID), zap.String("p", p.PostID), zap.Int8("direction", p.Direction))
	return redis.VoteForPost(c, strconv.FormatInt(userID, 10), p.PostID, float64(p.Direction))

}
