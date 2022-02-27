package router

import (
	"wj_rear/controller"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	r.POST("/api/user/register", controller.Register)

	r.Run()
}
