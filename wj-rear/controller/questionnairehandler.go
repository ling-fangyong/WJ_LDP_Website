package controller

import (
	"log"
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type QuesJson struct {
	WjID    uint `json:"WjId" binding:"required"`
	Options []struct {
		Title string `json:"title"`
		OpID  uint   `json:"opId"`
	} `json:"options"`
	QuesID uint   `json:"QuesId"`
	Type   int8   `json:"type"`
	Title  string `json:"title"`
}

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

func UpdateQuestion(ctx *gin.Context) {
	var quesfront QuesJson
	if err := ctx.BindJSON(&quesfront); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": "422",
			"msg":  "输入格式出错",
		})
		return
	}
	var ques model.Question
	var options []model.Option
	log.Println(quesfront)
	log.Println(quesfront.WjID, "   ", quesfront.QuesID)
	if quesfront.QuesID != 0 { //update question
		if err := database.DB.First(&ques, quesfront.QuesID).Error; err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "问题不存在",
			})
			return
		}
		//更新问题
		ques.Title = quesfront.Title
		ques.QuesType = quesfront.Type
		database.DB.Save(&ques)

		hash := make(map[uint]bool)
		//hash存储更新后仍存在的选项
		for _, option := range quesfront.Options {
			hash[option.OpID] = true
		}

		if err := database.DB.Where("Question_Id=?", quesfront.QuesID).Find(&options).Error; err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "选项查询失败",
			})
			return
		}
		log.Println(options)
		log.Println(hash)
		//删除不再存在的选项
		for _, option := range options {
			if _, ok := hash[option.ID]; !ok {
				database.DB.Delete(&model.Option{}, option.ID)
			}
		}

		//更新选项 增加选项
		for _, option := range quesfront.Options {
			var op model.Option
			op.QuestionId = quesfront.QuesID
			op.Title = option.Title
			if option.OpID != 0 {
				op.ID = option.OpID
				database.DB.Save(&op) //更新
			} else {
				database.DB.Create(&op) //增加
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新问题成功",
		})

	} else {
		ques.WjId = quesfront.WjID
		ques.Title = quesfront.Title
		ques.QuesType = quesfront.Type
		database.DB.Create(&ques)
		for _, option := range quesfront.Options {
			var op model.Option
			op.QuestionId = ques.ID
			op.Title = option.Title
			database.DB.Create(&op)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "创建问题成功",
		})
	}

}
