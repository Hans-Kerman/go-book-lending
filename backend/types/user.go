package types

import "github.com/Hans-Kerman/go-book-lending/backend/models"

type NewUser struct {
	UserName string `binding:"required,min=2,max=50" json:"user_name"`
	Password string `binding:"required,min=6" json:"password"`
}

type CtxEnvUser struct {
	ID       uint
	UserName string
	Role     models.UserRole
}
