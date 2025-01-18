package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// 获取当前登录用户userID
func getCurrentUser(c *gin.Context) (int64, error) {
	user, ok := c.Get(CtxUserIDKey)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	uid, ok := user.(int64)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	return uid, nil
}

func getPageInfo(c *gin.Context) (offset, limit int64) {
	offsetStr := c.Query("page")
	limitStr := c.Query("page_size")

	var (
		err error
	)

	offset, err = strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 1
	}

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	return
}
