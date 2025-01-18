package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	"web_app/models"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务调用

const secret = "cp"

var (
	ErrorUserNotExist        = errors.New("用户不存在")
	ErrorUserExist           = errors.New("用户已存在")
	ErrorUserInvalidPassword = errors.New("密码错误")
)

func CheckUserExist(username string) (err error) {
	// 根据用户名查找用户
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}
func GetUserByID(userID int64) (user *models.User, err error) {
	// 根据用户ID查找用户
	user = new(models.User)
	sqlStr := `select username from user where user_id = ?`
	if err = db.Get(user, sqlStr, userID); err != nil {
		return user, err
	}
	return
}

func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行sql语句入库
	sqlStr := `insert into user (user_id, username, password) values (?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		return err
	}
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登陆的密码
	sqlStr := `select user_id, username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}

	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorUserInvalidPassword
	}
	return
}
