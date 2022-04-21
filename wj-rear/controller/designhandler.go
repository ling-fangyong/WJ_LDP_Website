package controller

import (
	"fmt"
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
	DataMin float64      `json:"datamin"` //连续型数据最大值
	DataMax float64      `json:"datamax"` //连续型数据最小值
}

type OptionJson struct {
	Title  string `json:"title"`
	OpID   uint   `json:"opId"`
	CalcOp int    `json:"CalcOp"`
}
type DeleteJson struct {
	WjID   uint `json:"WjId"`
	QuesId uint `json:"QuesId"`
}

type QuesRetJson struct {
	WjID          uint         `json:"WjId" binding:"required"`
	Options       []OptionJson `json:"options"`
	QuesID        uint         `json:"QuesId"`
	Type          int8         `json:"type"`
	Title         string       `json:"title"`
	RadioValue    int          `json:"radiovalue"`
	CheckboxValue []int        `json:"checkboxValue"`
	Textvalue     string       `json:"textValue"`
	DataMin       float64      `json:"DataMin"` //连续型数据最大值
	DataMax       float64      `json:"DataMax"` //连续型数据最小值
}

func UpdateQuestionaire(ctx *gin.Context) {
	// title := ctx.PostForm("title")
	// desc := ctx.PostForm("desc")
	// var quetionaure model.Questionnaire
	// var WjId string
	// if err := ctx.BindJSON(&quetionaure); err != nil {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"code": 422,
	// 		"msg":  "输入格式不符合要求",
	// 	})
	// 	return
	// }
	// if WjId = ctx.PostForm("ID"); WjId != "" {
	// 	if err := database.DB.First(&quetionaure, WjId).Error; err != nil {
	// 		ctx.JSON(http.StatusOK, gin.H{
	// 			"code": 422,
	// 			"msg":  "问卷不存在",
	// 		})
	// 		return
	// 	}
	// } else {
	// 	session := sessions.Default(ctx)
	// 	if uid := session.Get("user_id"); uid != nil {
	// 		quetionaure.UserId = uid.(uint)
	// 	} //不做进一步判断，该一系列操作将放在auth鉴权下
	// }

	// var questionaire model.Questionnaire
	// if err := ctx.BindJSON(&questionaire); err != nil {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"code": 422,
	// 		"msg":  "输入格式不符合要求",
	// 	})
	// 	return
	// }
	// if questionaire.ID == 0 {
	// 	session := sessions.Default(ctx)
	// 	if uid := session.Get("user_id"); uid != nil {
	// 		questionaire.UserId = uid.(uint)
	// 	} //不做进一步判断，该一系列操作将放在auth鉴权下
	// }
	// fmt.Println("title:", questionaire.Title)
	// fmt.Println("desc:", questionaire.Desc)
	// if questionaire.ID != 0 {
	// 	database.DB.Save(&questionaire)
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"code": 200,
	// 		"msg":  "更新问卷成功",
	// 	})
	// } else {
	// 	database.DB.Create(&questionaire)
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"code": 200,
	// 		"msg":  "创建问卷成功",
	// 	})
	// }

	var questionairejson QuestionaireJson
	var questionaire model.Questionnaire
	if err := ctx.BindJSON(&questionairejson); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "输入格式不符合要求",
		})
		return
	}
	if questionairejson.WjId == 0 {
		session := sessions.Default(ctx)
		if uid := session.Get("user_id"); uid != nil {
			questionaire.UserId = uid.(uint)
		} //不做进一步判断，该一系列操作将放在auth鉴权下
	} else {
		if err := database.DB.First(&questionaire, questionairejson.WjId).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 422,
				"msg":  "问卷不存在",
			})
			return
		}
	}
	// fmt.Println("title:", questionaire.Title)
	// fmt.Println("desc:", questionaire.Desc)
	questionaire.Title = questionairejson.Title
	questionaire.Desc = questionairejson.Desc
	if questionairejson.WjId != 0 {
		database.DB.Save(&questionaire)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新问卷成功",
		})
	} else {
		database.DB.Create(&questionaire)
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
		ctx.JSON(http.StatusOK, gin.H{
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
			ctx.JSON(http.StatusOK, gin.H{
				"code": 422,
				"msg":  "问题不存在",
			})
			return
		}
		//更新问题
		ques.Title = quesfront.Title
		ques.QuesType = quesfront.Type
		ques.DataMin = quesfront.DataMin
		ques.DataMax = quesfront.DataMax
		database.DB.Save(&ques)

		hash := make(map[uint]bool)
		//hash存储更新后仍存在的选项
		for _, option := range quesfront.Options {
			hash[option.OpID] = true
		}

		if err := database.DB.Where("Question_Id=?", quesfront.QuesID).Find(&options).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{
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
		ques.DataMin = quesfront.DataMin
		ques.DataMax = quesfront.DataMax
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
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "输入格式不符合要求",
		})
		return
	}
	if DeleteArg.WjID != 0 {
		tx := database.DB.Begin()
		if err := database.DB.Where("Id=?", DeleteArg.WjID).Delete(&model.Questionnaire{}).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 422,
				"msg":  "问卷删除失败",
			})
			tx.Rollback()
			return
		}
		var ques []model.Question
		if err := database.DB.Where("Wj_Id=?", DeleteArg.WjID).Find(&ques).Error; err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 422,
				"msg":  "问题查询失败",
			})
			tx.Rollback()
			return
		} else {
			log.Println(ques)
			for _, Ques := range ques {
				if ok, msg := DeleteQuestion(Ques.ID); !ok {
					ctx.JSON(http.StatusOK, msg)
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
			_, msg := DeleteQuestion(DeleteArg.QuesId)
			ctx.JSON(http.StatusOK, msg)
			return
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
		ctx.JSON(http.StatusOK, gin.H{
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
		// if err := ctx.BindJSON(&WjId); err != nil {
		// 	ctx.JSON(http.StatusOK, gin.H{
		// 		"code": 422,
		// 		"msg":  "输入格式不符合要求",
		// 	})
		// 	return
		// }
		WjId = ctx.PostForm("WjId")
		fmt.Println(WjId)
	}
	var ques []model.Question
	if err := database.DB.Where("Wj_Id=?", WjId).Find(&ques).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "问卷问题获取失败",
		})
		return
	} else {
		var question []QuesRetJson
		for _, item := range ques {
			var quesItem QuesRetJson
			quesItem.WjID = item.WjId
			quesItem.QuesID = item.ID
			quesItem.Title = item.Title
			quesItem.Type = item.QuesType
			quesItem.DataMax = item.DataMax
			quesItem.DataMin = item.DataMin
			if quesItem.Type == 1 {
				quesItem.RadioValue = 0
			} else if quesItem.Type == 2 {
				quesItem.CheckboxValue = make([]int, 0)
			}
			var options []model.Option
			if err := database.DB.Where("question_id=?", item.ID).Find(&options).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
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
