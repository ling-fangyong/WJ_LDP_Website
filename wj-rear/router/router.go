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

			authed.POST("/design/UpdateQuestionaire", controller.UpdateQuestionaire)

			authed.POST("/design/UpdateQuestion", controller.UpdateQuestion)

			authed.DELETE("/design/DeleteQuestionaire", controller.DeleteQuestionaire)

			authed.GET("/show/ShowQuestionaires", controller.ShowQuestionaires)

			authed.POST("/show/ShowQuestions", controller.ShowQuestions)
		}

		answer := api.Group("/answer")
		{
			answer.GET("/GetQuestionaire", controller.GetQuestionaire)

			answer.POST("/SubmitQues", controller.SubmitQues)
		}
	}

	r.Run(":3030")
}
