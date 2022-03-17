package controller

import (
	"log"
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type QuestionaireJson struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	WjId  uint   `json:"wjid"`
}

type QuesJson struct {
	WjID    uint         `json:"WjId" binding:"required"`
	Options []OptionJson `json:"options"`
	QuesID  uint         `json:"QuesId"`
	Type    int8         `json:"type"`
	Title   string       `json:"title"`
}

type OptionJson struct {
	Title string `json:"title"`
	OpID  uint   `json:"opId"`
}
type DeleteJson struct {
	WjID   uint `json:"WjID"`
	QuesId uint `json:"QuesId"`
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

//删除问卷或问题，通过复用删除问题来完成
func DeleteQuestionaire(ctx *gin.Context) {
	var DeleteArg DeleteJson
	if err := ctx.BindJSON(&DeleteArg); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "输入格式不符合要求",
		})
		return
	}
	if DeleteArg.WjID != 0 {
		tx := database.DB.Begin()
		if err := database.DB.Where("Id=?", DeleteArg.WjID).Delete(&model.Questionnaire{}).Error; err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "问卷删除失败",
			})
			tx.Rollback()
			return
		}
		var ques []model.Question
		if err := database.DB.Where("Wj_Id=?", DeleteArg.WjID).Find(&ques).Error; err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "问题查询失败",
			})
			tx.Rollback()
			return
		} else {
			log.Println(ques)
			for _, Ques := range ques {
				if ok, msg := DeleteQuestion(Ques.ID); !ok {
					ctx.JSON(http.StatusUnprocessableEntity, msg)
					tx.Rollback()
					return
				}
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "删除问卷成功",
			})
		}
		tx.Commit()
	} else {
		if DeleteArg.QuesId != 0 {
			if ok, msg := DeleteQuestion(DeleteArg.QuesId); !ok {
				ctx.JSON(http.StatusUnprocessableEntity, msg)
				return
			}
		}
	}
}

//删除问题
func DeleteQuestion(QuesId uint) (bool, gin.H) {
	var ques model.Question
	if err := database.DB.First(&ques, QuesId).Error; err != nil {
		return false, gin.H{
			"code": 422,
			"msg":  "问题不存在",
		}
	} else {
		tx := database.DB.Begin()
		//删除问题
		database.DB.Delete(&model.Question{}, QuesId)

		if err := database.DB.Where("Question_id=?", QuesId).Delete(&model.Option{}).Error; err != nil {
			tx.Rollback()
			return false, gin.H{
				"code": 422,
				"msg":  "删除选项失败",
			}
		}

		if err := database.DB.Where("Question_id=?", QuesId).Delete(&model.Answer{}).Error; err != nil {
			tx.Rollback()
			return false, gin.H{
				"code": 422,
				"msg":  "删除答案失败",
			}
		}
		tx.Commit()
		return true, gin.H{
			"code": 200,
			"msg":  "删除问题成功",
		}
	}
}

//获取该用户问卷列表
func ShowQuestionaires(ctx *gin.Context) {
	session := sessions.Default(ctx)
	uid := session.Get("user_id")
	var questionaire []model.Questionnaire
	if err := database.DB.Where("user_id=?", uid).Find(&questionaire).Error; err != nil {
		log.Println(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "问卷列表获取失败",
		})
		return
	} else {
		var quesjson []QuestionaireJson
		for _, item := range questionaire {
			quesjson = append(quesjson, QuestionaireJson{
				Title: item.Title,
				Desc:  item.Desc,
				WjId:  item.ID,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取问卷列表成功",
			"data": quesjson,
		})
	}
}

//获取问题列表
func ShowQuestions(ctx *gin.Context) {
	WjId := ctx.Query("id")
	if WjId == "" {
		WjId = ctx.PostForm("WjId")
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
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取问题列表成功",
			"data": question,
		})
	}
}
