package controller

import (
	"github.com/gin-gonic/gin"

	"web_app/models"
)

func PostVoteHandler(c *gin.Context) {
	// 获取参数以及参数校验
	p := new(models.VoteData)
	c.ShouldBindJSON(p)
}
