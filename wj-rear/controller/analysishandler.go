package controller

import (
	"net/http"
	"strconv"
	"wj_rear/database"
	"wj_rear/model"

	"wj_rear/algorithm"

	"github.com/gin-gonic/gin"
)

const epsilon = 1

type result struct {
	ans_int int
	ansCalc int
}

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
			} else if quesItem.Type == 2 {
				if err := database.DB.Raw("SELECT a.id AS OpID,a.title,a.DummyValueCnt,b.CalcOp FROM ( SELECT id,title,DummyValueCnt FROM `options` WHERE question_id = ?) AS a LEFT JOIN ( SELECT op_id,COUNT(op_id) AS CalcOp FROM answers WHERE question_id = ? GROUP BY op_id) As b ON a.id = b.op_id;", quesItem.QuesID, quesItem.QuesID).Find(&quesItem.Options).Error; err != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"code": 422,
						"msg":  "问题选项获取失败",
					})
					return
				} else {
					// fmt.Println(quesItem.Options)

					res := make([]int, 2*len(quesItem.Options))
					for index, option := range quesItem.Options {
						res[index] = option.CalcOp
					}
					for index := range quesItem.Options {
						res[index+len(quesItem.Options)] = int(quesItem.Options[index].DummyValueCnt)
					}
					// fmt.Println("---------------------")
					// fmt.Println(res)
					// fmt.Println(algorithm.GRR(res, len(res), epsilon))
					copy(res, algorithm.GRR(res, len(res), epsilon))

					// fmt.Println("---------------------")
					// fmt.Println(res)
					// fmt.Println("---------------------")
					for index := range quesItem.Options {
						quesItem.Options[index].CalcOp = res[index] * len(quesItem.Options)
					}
					// fmt.Println(quesItem.Options)
					question = append(question, quesItem)
				}
			} else if quesItem.Type == 3 {
				quesItem.DataMax = item.DataMax
				quesItem.DataMin = item.DataMin
				var res []result
				if err := database.DB.Raw("SELECT ans_int,COUNT(ans_int) AS num FROM `answers` WHERE question_id = ? GROUP BY ans_int", quesItem.QuesID).Find(&res).Error; err != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"code": 422,
						"msg":  "问题答案获取失败",
					})
					return
				} else {
					if len(res) != 0 {
						var tem = make([]int, 2)
						tem[0] = res[0].ansCalc
						if len(res) == 2 {
							tem[1] = res[1].ansCalc
						} else {
							tem[1] = 0
						}
						ans_int0 := res[0].ans_int
						var ans_int1 int
						if ans_int0 == 1 {
							ans_int1 = -1
						} else {
							ans_int1 = 1
						}
						copy(tem, algorithm.GRR(tem, 2, epsilon))
						value := (float64(ans_int0*tem[0])+float64(ans_int1*tem[1]))/float64(tem[0]+tem[1])*((item.DataMax-item.DataMin)/2) + (item.DataMax+item.DataMin)/2
						quesItem.Textvalue = strconv.FormatFloat(value, 'E', -1, 64)
					} else {
						quesItem.Textvalue = "不存在填写值"
					}
					question = append(question, quesItem)
				}
			}
		}
		//fmt.Println(question)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取问题列表成功",
			"data": question,
		})
	}
}
