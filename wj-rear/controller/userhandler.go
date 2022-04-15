package controller

import (
	"log"
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "用户名或密码不符合要求",
		})
	} else {
		var count int64 = 0
		database.DB.Model(&model.User{}).Where("Name=?", user.Name).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 422,
				"msg":  "用户名已注册",
			})
		} else {
			if err := user.HashPassword(); err != nil {
				log.Println(err)
				ctx.JSON(http.StatusServiceUnavailable, gin.H{
					"code": 422,
					"msg":  "密码hash错误",
				})
			}
			if err := database.DB.Create(&user).Error; err != nil {
				log.Println(err)
				ctx.JSON(http.StatusServiceUnavailable, gin.H{
					"code": 422,
					"msg":  "数据创建失败",
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

func Login(ctx *gin.Context) {
	var user model.User //前端传输过来的用户信息
	if err := ctx.ShouldBind(&user); err != nil {
		log.Println(err)
		log.Println(user)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "用户名或密码不符合要求",
		})
	} else {
		var dbUser model.User //数据库中的用户信息
		if err := database.DB.Model(&model.User{}).Where("Name=?", user.Name).First(&dbUser).Error; err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 422,
				"msg":  "用户名或者密码错误",
			})
		} else {
			if err := dbUser.CheckPassword(user.Password); err != nil {
				log.Println(err)
				ctx.JSON(http.StatusOK, gin.H{
					"code": 422,
					"msg":  "用户名或者密码错误",
				})
			} else {
				session := sessions.Default(ctx)
				session.Clear()
				session.Set("user_id", dbUser.ID)
				session.Save()
				ctx.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "登陆成功",
				})
			}
		}

	}
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})
}

func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "pong",
	})
}
