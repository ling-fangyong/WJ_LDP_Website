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

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/ping", controller.Ping)

			user.POST("/register", controller.Register)

			user.POST("/login", controller.Login)

			user.DELETE("/logout", controller.Logout)
		}
	}

	r.Run(":8080")
}
