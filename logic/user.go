package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

// 存放业务逻辑代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户存不存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成uid
	userID := snowflake.GetID()

	// 构造一个User实例
	u := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 3. 保存进数据库
	err = mysql.InsertUser(u)
	if err != nil {
		return err
	}

	return nil
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 传递的是指针，就能拿到user.userID
	if err = mysql.Login(user); err != nil {
		// 登录失败
		return nil, err
	}

	// 生成JWT
	aToken, rToken, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	user.AToken = aToken
	user.RToken = rToken
	return
}
