package middleware

import (
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		if uid := session.Get("user_id"); uid != nil {
			user, err := database.GetUserById(uid)
			if err == nil {
				ctx.Set("user", &user)
			}
		}
		ctx.Next() //是否一定需要加，如果不加，正常执行后续handler好像也没问题
	}
}
func AuthUserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if user, _ := ctx.Get("user"); user != nil {
			if _, ok := user.(*model.User); ok {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "需要登录",
		})
		ctx.Abort()
	}
}
