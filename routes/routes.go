package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"web_app/controller"
	"web_app/logger"
	"web_app/middleware"
	"web_app/pkg/snowflake"
)

func SetUp() *gin.Engine {
	r := gin.New()
	// 重新使用zap新写中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, cast.ToString(snowflake.GetID()))
	})

	v1 := r.Group("/api/v1")

	// 注册业务路由
	v1.POST("/signUp", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostHandler)
		v1.GET("/post", controller.GetPostListHandler)

		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
