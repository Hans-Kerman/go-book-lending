package types

type NewUser struct {
	UserName string `binding:"required,min=2,max=50" json:"user_name"`
	Password string `binding:"required,min=6" json:"password"`
}
