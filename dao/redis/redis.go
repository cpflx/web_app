package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"web_app/settings"
)

var rdb *redis.Client

func Init(conf *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Pass,     // 密码
		DB:       conf.DB,       // 数据库
		PoolSize: conf.PoolSize, // 连接池大小
	})

	ctx := context.Background()
	if _, err = rdb.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func Close() {
	_ = rdb.Close()
}
