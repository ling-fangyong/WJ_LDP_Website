package controller

import (
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UpdateQuestionaire(ctx *gin.Context) {
	title := ctx.PostForm("title")
	desc := ctx.PostForm("desc")
	var quetionaure model.Questionnaire
	var WjId string
	if WjId = ctx.PostForm("ID"); WjId != "" {
		if err := database.DB.First(&quetionaure, WjId).Error; err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "问卷不存在",
			})
			return
		}
	} else {
		session := sessions.Default(ctx)
		if uid := session.Get("user_id"); uid != nil {
			quetionaure.UserId = uid.(uint)
		} //不做进一步判断，该一系列操作将放在auth鉴权下
	}
	quetionaure.Title = title
	quetionaure.Desc = desc
	if WjId != "" {
		database.DB.Save(&quetionaure)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新问卷成功",
		})
	} else {
		database.DB.Create(&quetionaure)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建问卷成功",
		})
	}

}
