package controller

import (
	"log"
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  err,
		})
	} else {
		var count int64 = 0
		database.DB.Model(&model.User{}).Where("Name=?", user.Name).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户名已注册",
			})
		} else {
			if err := database.DB.Create(&user).Error; err != nil {
				log.Println(err)
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"code": 422,
					"msg":  err,
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "注册成功",
				})
			}
		}
	}
}
