package redis

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math"
	"time"
)

// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm

// 本项目使用简化版的投票分数
// 投一票加432分， 86400/200-》200票帖子就能续一天

/*
投票的几种情况：
direction = 1 为赞成票:2种情况
	1：用户第一次投赞成票 差值绝对值1 +432
	2：之前反对，现在赞成 差值绝对值2 +432*2

direction = 0 2种情况
	1：之前投赞成票，现在取消赞成票 差值绝对值1 -432
	2：之前反对，现在取消反对票 差值绝对值1 +432

direction = -1 为反对票:2种情况
	1: 用户第一次投反对票 差值绝对值1 -432
	2: 之前赞成，现在反对票 差值绝对值2

投票的限制：
每个帖子自发表之日一周内才可投票，超过则无法投票
 	1. 到期之后将redis 中该帖子的投票数据存储到mysql中
	2. 删除redis中该帖子的投票数据 KeyPostVotedZSetPF

*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 分数每次+432
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(ctx *gin.Context, postID int64) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	if err != nil {
		return err
	}
	_, err = pipeline.Exec(ctx)
	return err
}

func VoteForPost(ctx *gin.Context, userID, postID string, value float64) (err error) {
	// 1. 判断投票限制
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if int(float64(time.Now().Unix())-postTime) > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}

	// 2. 更新帖子的分数
	// 查当前用户给当前帖子之前的投票记录
	ov := rdb.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF)+postID, userID).Val()
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value)

	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF)+postID, userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF)+postID, redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)

	return err
}
