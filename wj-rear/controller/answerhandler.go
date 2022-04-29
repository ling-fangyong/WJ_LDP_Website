package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"wj_rear/database"
	"wj_rear/model"

	"github.com/gin-gonic/gin"
)

type GetQuesJson struct {
	Questionaire QuestionaireJson `json:"Questionaire"`
	Ques         []QuesRetJson    `json:"Ques"`
}

type AnsJson struct {
	WjId      uint          `json:"WjId"`
	QuesAndOp []QuesRetJson `json:"QuesAndOp"`
}

// type AnsJson struct {
// 	WjId uint `json:"WjId"`
// 	Ans  []struct {
// 		QuesId    uint   `json:"QuesID"`
// 		AnsInt    uint   `json:"AnsInt"`
// 		AnsString string `json:"AnsString"`
// 		QuesType  int8   `json:"QuesType"`
// 		//TODO:暂时完成选择题，填空题只允许数值型数据进行分析，此处还需要存储值范围
// 	} `json:"ans"`
// }

func GetQuestionaire(ctx *gin.Context) {
	WjId := ctx.PostForm("WjId")
	// fmt.Println(WjId)
	var QuestionaireModel model.Questionnaire
	if err := database.DB.First(&QuestionaireModel, WjId).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "问卷获取失败",
		})
		return
	}
	var ques []model.Question
	fmt.Println(WjId)

	if err := database.DB.Where("Wj_Id=?", WjId).Find(&ques).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "问卷问题获取失败",
		})
		return
	} else {
		var question []QuesRetJson
		// fmt.Println("ques")
		fmt.Println(ques)
		for _, item := range ques {
			var quesItem QuesRetJson
			quesItem.QuesID = item.ID
			quesItem.Title = item.Title
			quesItem.Type = item.QuesType
			quesItem.WjID = item.WjId
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
		var QuesShow GetQuesJson
		QuesShow.Questionaire.WjId = QuestionaireModel.ID
		QuesShow.Questionaire.Title = QuestionaireModel.Title
		QuesShow.Questionaire.Desc = QuestionaireModel.Desc
		QuesShow.Ques = question
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取问题列表成功",
			"data": QuesShow,
		})
	}
}

// func SubmitQues(ctx *gin.Context) {
// 	var ansJson AnsJson
// 	if err := ctx.BindJSON(&ansJson); err != nil {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"code": "422",
// 			"msg":  "输入格式出错",
// 		})
// 		return
// 	}

// 	//判断问卷是否存在
// 	var count int64
// 	if err := database.DB.Model(&model.Questionnaire{}).Where("Id=?", ansJson.WjId).Count(&count).Error; err != nil || count < 1 {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"code": "422",
// 			"msg":  "问卷不存在或查询失败",
// 		})
// 		return
// 	}
// 	log.Println("Ans  ", ansJson)
// 	tx := database.DB.Begin()
// 	for _, item := range ansJson.Ans {
// 		var ans model.Answer
// 		ans.WjId = ansJson.WjId
// 		ans.QuesType = item.QuesType
// 		ans.QuestionId = item.QuesId
// 		if item.QuesType == 1 || item.QuesType == 2 { //TODO:当类型为2时还要增加两个变量代表数值范围，暂时先完成选择题测试
// 			ans.AnsInt = item.AnsInt
// 			if err := database.DB.Create(&ans).Error; err != nil {
// 				ctx.JSON(http.StatusOK, gin.H{
// 					"code": "422",
// 					"msg":  "答案创建失败",
// 				})
// 				tx.Rollback()
// 				return
// 			}
// 		} else if item.QuesType == 3 {
// 			ans.AnsString = item.AnsString

// 			if err := database.DB.Create(&ans).Error; err != nil {
// 				ctx.JSON(http.StatusOK, gin.H{
// 					"code": "422",
// 					"msg":  "答案创建失败",
// 				})
// 				tx.Rollback()
// 				return
// 			}
// 		}
// 	}
// 	tx.Commit()
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"code": 200,
// 		"msg":  "问卷提交成功",
// 	})
// }

func SubmitQues(ctx *gin.Context) {
	var ansJson AnsJson
	if err := ctx.BindJSON(&ansJson); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "422",
			"msg":  "输入格式出错",
		})
		return
	}
	var count int64
	if err := database.DB.Model(&model.Questionnaire{}).Where("Id=?", ansJson.WjId).Count(&count).Error; err != nil || count < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "422",
			"msg":  "问卷不存在或查询失败",
		})
		return
	}

	tx := database.DB.Begin()
	for _, item := range ansJson.QuesAndOp {
		if item.Type == 1 {
			var ans model.Answer
			ans.QuestionId = item.QuesID
			ans.OpId = item.Options[item.RadioValue].OpID
			fmt.Println(ans)
			if err := database.DB.Create(&ans).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": "422",
					"msg":  "答案创建失败",
				})
				tx.Rollback()
				return
			}
		} else if item.Type == 2 {
			// //理论上checkValue一定只有一个值，但是懒得修改了
			// for _, checkValue := range item.CheckboxValue {
			// 	var ans model.Answer
			// 	//当出现超过选项数值选项值时跳过
			// 	if checkValue >= len(item.Options) {
			// 		continue
			// 	}
			// 	ans.QuestionId = item.QuesID
			// 	ans.OpId = item.Options[checkValue].OpID
			// 	fmt.Println(ans)
			// 	if err := database.DB.Create(&ans).Error; err != nil {
			// 		ctx.JSON(http.StatusOK, gin.H{
			// 			"code": "422",
			// 			"msg":  "答案创建失败",
			// 		})
			// 		tx.Rollback()
			// 		return
			// 	}
			// }
			for _, checkValue := range item.CheckboxValue {
				var ans model.Answer
				//当出现超过选项数值选项值时跳过
				if checkValue >= len(item.Options) {
					temOpId := item.Options[checkValue-len(item.Options)].OpID
					var temOp model.Option
					if err := database.DB.Find(&temOp, temOpId).Error; err != nil {
						ctx.JSON(http.StatusOK, gin.H{
							"code": "422",
							"msg":  "查询失败",
						})
						tx.Rollback()
						return
					}
					temOp.DummyValueCnt++
					database.DB.Save(&temOp)
				} else {
					ans.QuestionId = item.QuesID
					ans.OpId = item.Options[checkValue].OpID
					fmt.Println(ans)
					if err := database.DB.Create(&ans).Error; err != nil {
						ctx.JSON(http.StatusOK, gin.H{
							"code": "422",
							"msg":  "答案创建失败",
						})
						tx.Rollback()
						return
					}
				}
			}
		} else if item.Type == 3 {
			var ans model.Answer
			ans.QuestionId = item.QuesID
			ans.AnsInt, _ = strconv.Atoi(item.Textvalue)
			if err := database.DB.Create(&ans).Error; err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": "422",
					"msg":  "答案创建失败",
				})
				tx.Rollback()
				return
			}
		}
		//TODO:type=3的连续型数据以及键值型待做，改为集合型，由多选题实现
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问卷提交成功",
	})
}
