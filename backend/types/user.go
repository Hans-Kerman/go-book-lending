package types

import (
	"gorm.io/gorm"
)

type NewUser struct {
	gorm.Model
	UserName string `binding:"required" json:"user_name"`
	Password string `binding:"required" json:"password"`
}
