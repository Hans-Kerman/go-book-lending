package models

import (
	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string         `gorm:"type:varchar(50);not null;uniqueIndex"`
	Role     types.UserRole `gorm:"not null"` //角色只有Admin(可以添加书籍、查询所有人的数据)和Reader(可以借阅、查询自己数据)
	Password string         `gorm:"not null"`
}
