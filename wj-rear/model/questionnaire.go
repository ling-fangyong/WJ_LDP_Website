package model

import "gorm.io/gorm"

//问卷表
type Questionnaire struct {
	gorm.Model      //自带id为问卷ID
	UserId     uint //用户ID
	Title      string
	Desc       string
}

//问题表
type Question struct {
	gorm.Model        //自带id为问题ID
	WjId       uint   //问卷ID
	Title      string //问题
	QuesType   int8   //1代表单选题 2代表多选题 3代表数值填写  前两者1、2可以对应离散型数据 3代表连续型数据需要数据范围
	// OptionNum  int8   //选项个数
	DataMin float64 //连续型数据最大值
	DataMax float64 //连续型数据最小值
}

//选项表
type Option struct {
	gorm.Model
	QuestionId    uint   //问题ID
	Title         string //选项
	DummyValueCnt uint   //存储着多选时的虚假值
}

//回答表
// type Answer struct {
// 	gorm.Model
// 	WjId       uint //问卷ID
// 	QuestionId uint //问题ID
// 	OpId       uint //选项ID
// 	QuesType   int8 //问题类型暂时不加
// 	AnsInt     uint //TODO:问题类型为1时存储选择选项的ID 问题类型为2时存储映射并转换后的0 1需要添加存储范围  3可能用于键值型数据 需要补充
// 	AnsString  string
// }
type Answer struct {
	gorm.Model
	QuestionId uint //获取统计答案时，问题类型直接决定答案类型
	OpId       uint //选择题填写该Id
	AnsInt     int  //填空题填写该答案 由于会经过变换值只为-1 1
	//键值型再说
}
