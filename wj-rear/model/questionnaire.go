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
	QuesType   int8   //0代表单选题 1代表多选题 2代表数值填写  前两者0、1可以对应离散型数据 2代表连续型数据需要数据范围
	OptionNum  int8   //选项个数
}

//选项表
type Option struct {
	gorm.Model
	QuestionId uint   //问题ID
	Title      string //选项
}

//回答表
type Answer struct {
	gorm.Model
	WjId       uint //问卷ID
	QuestionId uint //问题ID
	QuesType   int8 //问题类型暂时不加
	AnsInt     uint //问题类型为0时存储选择选项的ID   问题类型为1时存储映射并转换后的0 1
	AnsString  string
}
