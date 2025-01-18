package controller

import (
	"github.com/gin-gonic/gin"

	"web_app/models"
)

func PostVoteHandler(c *gin.Context) {
	// 获取参数以及参数校验
	p := new(models.ParamVoteData)
	err := c.ShouldBindJSON(p)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
}
