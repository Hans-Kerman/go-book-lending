package types

type UserRole uint8

const (
	Admin UserRole = iota
	Reader
)

type NewUser struct {
	UserName string `binding:"required,min=2,max=50" json:"user_name"`
	Password string `binding:"required,min=6" json:"password"`
}

type CtxEnvUser struct {
	ID       uint
	UserName string
	Role     UserRole
}
