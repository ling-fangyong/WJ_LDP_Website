package controller

import (
	"log"
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-gonic/gin"
)

type GetQuesJson struct {
	Questionnaire QuestionaireJson
	Ques          []QuesJson
}

type AnsJson struct {
	WjId uint `json:"WjId"`
	Ans  []struct {
		QuesId    uint   `json:"QuesID"`
		AnsInt    uint   `json:"AnsInt"`
		AnsString string `json:"AnsString"`
		QuesType  int8   `json:"QuesType"`
	} `json:"ans"`
}

func GetQuestionaire(ctx *gin.Context) {
	WjId := ctx.PostForm("WjId")
	var QuestionaireModel model.Questionnaire
	if err := database.DB.First(&QuestionaireModel, WjId).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "问卷获取失败",
		})
		return
	}
	var ques []model.Question
	if err := database.DB.Where("Wj_Id=?", WjId).Find(&ques).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "问卷问题获取失败",
		})
		return
	} else {
		var question []QuesJson
		for _, item := range ques {
			var quesItem QuesJson
			quesItem.QuesID = item.ID
			quesItem.Title = item.Title
			quesItem.Type = item.QuesType
			var options []model.Option
			if err := database.DB.Where("question_id=?", item.ID).Find(&options).Error; err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"code": 422,
					"msg":  "问题选项获取失败",
				})
				return
			} else {
				for _, option := range options {
					quesItem.Options = append(quesItem.Options, OptionJson{
						OpID:  option.ID,
						Title: option.Title,
					})
				}

				question = append(question, quesItem)
			}
		}
		var QuesShow GetQuesJson
		QuesShow.Questionnaire.WjId = QuestionaireModel.ID
		QuesShow.Questionnaire.Title = QuestionaireModel.Title
		QuesShow.Questionnaire.Desc = QuestionaireModel.Desc
		QuesShow.Ques = question
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取问题列表成功",
			"data": QuesShow,
		})
	}
}

func SubmitQues(ctx *gin.Context) {
	var ansJson AnsJson
	if err := ctx.BindJSON(&ansJson); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": "422",
			"msg":  "输入格式出错",
		})
		return
	}

	//判断问卷是否存在
	var count int64
	if err := database.DB.Model(&model.Questionnaire{}).Where("Id=?", ansJson.WjId).Count(&count).Error; err != nil || count < 1 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": "422",
			"msg":  "问卷不存在或查询失败",
		})
		return
	}
	log.Println("Ans  ", ansJson)
	tx := database.DB.Begin()
	for _, item := range ansJson.Ans {
		var ans model.Answer
		ans.WjId = ansJson.WjId
		ans.QuesType = item.QuesType
		ans.QuestionId = item.QuesId
		if item.QuesType == 1 { //
			ans.AnsInt = item.AnsInt

			if err := database.DB.Create(&ans).Error; err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"code": "422",
					"msg":  "答案创建失败",
				})
				tx.Rollback()
				return
			}
		} else if item.QuesType == 2 {
			ans.AnsString = item.AnsString

			if err := database.DB.Create(&ans).Error; err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"code": "422",
					"msg":  "答案创建失败",
				})
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问卷提交成功",
	})
}