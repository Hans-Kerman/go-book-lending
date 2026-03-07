package pkg

import (
	"time"

	"github.com/Hans-Kerman/go-book-lending/backend/config"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   uint            `json:"user_id"`
	UserName string          `json:"user_name"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, username string, role models.UserRole) (string, error) {
	JWTConfig := config.AppConfig.JWT

	claim := &UserClaims{
		UserID:   userID,
		UserName: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTConfig.ExpireTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(JWTConfig.SecretStr)
}
