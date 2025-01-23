package redis

// redis key

const (
	KeyPrefix          = "webApp:"
	KeyPostTimeZSet    = "post:time"   // zset 用来存储帖子的发布时间
	KeyPostScoreZSet   = "post:score"  // zset 用来存储帖子的分数
	KeyPostVotedZSetPF = "post:voted:" // zset 用来记录用户对帖子的投票状态 1赞成 -1 参数是post_id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
