package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"web_app/settings"
)

var db *sqlx.DB

func Init(conf *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		conf.User,
		conf.Pass,
		conf.Host,
		conf.Port,
		conf.DB,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect mysqlDB failed, err:%v\n", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(conf.MaxOpen)
	db.SetMaxIdleConns(conf.MaxIdle)
	return nil
}

func Close() {
	_ = db.Close()
}
