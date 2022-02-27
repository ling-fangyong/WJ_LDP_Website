package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `form:"name" json:"name" binding:"required,min=4,max=16"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}
