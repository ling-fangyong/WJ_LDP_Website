package controller

import (
	"net/http"
	"wj_rear/database"
	"wj_rear/model"

	"wj_rear/algorithm"

	"github.com/gin-gonic/gin"
)

const epsilon = 1

func AnalysisData(ctx *gin.Context) {
	WjId := ctx.PostForm("WjId")
	var ques []model.Question
	if err := database.DB.Where("Wj_Id=?", WjId).Find(&ques).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 422,
			"msg":  "问卷问题获取失败",
		})
		return
	} else {
		// fmt.Println(WjId)
		// fmt.Println(ques)
		var question []QuesRetJson
		for _, item := range ques {
			var quesItem QuesRetJson
			quesItem.WjID = item.WjId
			quesItem.QuesID = item.ID
			quesItem.Title = item.Title
			quesItem.Type = item.QuesType

			// fmt.Println(quesItem.QuesID)
			// fmt.Println(quesItem.Options)
			if quesItem.Type == 1 || quesItem.Type == 2 {
				if err := database.DB.Raw("SELECT a.id AS OpID,a.title,b.CalcOp FROM ( SELECT id,title FROM `options` WHERE question_id = ?) AS a LEFT JOIN ( SELECT op_id,COUNT(op_id) AS CalcOp FROM answers WHERE question_id = ? GROUP BY op_id) As b ON a.id = b.op_id;", quesItem.QuesID, quesItem.QuesID).Find(&quesItem.Options).Error; err != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"code": 422,
						"msg":  "问题选项获取失败",
					})
					return
				} else {
					// fmt.Println(quesItem.Options)
					res := make([]int, len(quesItem.Options))
					for index, option := range quesItem.Options {
						res[index] = option.CalcOp
					}
					// fmt.Println("---------------------")
					// fmt.Println(res)
					// fmt.Println(algorithm.GRR(res, len(res), epsilon))
					copy(res, algorithm.GRR(res, len(res), epsilon))
					// fmt.Println("---------------------")
					// fmt.Println(res)
					// fmt.Println("---------------------")
					for index := range quesItem.Options {
						quesItem.Options[index].CalcOp = res[index]
					}
					// fmt.Println(quesItem.Options)
					question = append(question, quesItem)
				}
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取问题列表成功",
			"data": question,
		})
	}
}
