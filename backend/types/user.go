package types

type UserRole uint8

const (
	Admin UserRole = iota
	Reader
)

type NewUser struct {
	UserName string `binding:"required,min=2,max=50" json:"user_name" example:"testuser"`
	Password string `binding:"required,min=6" json:"password" example:"password123"`
}

type CtxEnvUser struct {
	ID       uint
	UserName string
	Role     UserRole
}
