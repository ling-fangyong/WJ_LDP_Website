package database

import "wj_rear/model"

func migration() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Questionnaire{})
	DB.AutoMigrate(&model.Question{})
	DB.AutoMigrate(&model.Option{})
	DB.AutoMigrate(&model.Answer{})
}
