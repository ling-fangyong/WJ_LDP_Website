package router

import (
	"wj_rear/controller"
	"wj_rear/middleware"

	"github.com/gin-gonic/gin"
)

var secretKey = "wj_secret_key"

func Run() {
	r := gin.Default()

	r.Use(middleware.Session(secretKey))

	r.Use(middleware.CurrentUser())

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/ping", controller.Ping)

			user.POST("/register", controller.Register)

			user.POST("/login", controller.Login)

		}

		authed := api.Group("/")
		authed.Use(middleware.AuthUserLogin()) //登录验证
		{
			authed.DELETE("/user/logout", controller.Logout)

			authed.POST("/design/questionaire", controller.UpdateQuestionaire)
		}

	}

	r.Run(":8080")
}
