package database

import "wj_rear/model"

func migration() {
	DB.AutoMigrate(&model.User{})
}
