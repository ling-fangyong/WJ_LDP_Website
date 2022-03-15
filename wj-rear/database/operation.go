package database

import (
	"wj_rear/model"
)

func GetUserById(ID interface{}) (model.User, error) {
	var user model.User
	res := DB.First(&user, ID)
	return user, res.Error
}
