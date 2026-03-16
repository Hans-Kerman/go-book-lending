package pkg

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Hans-Kerman/go-book-lending/backend/config"
	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   uint           `json:"user_id"`
	UserName string         `json:"user_name"`
	Role     types.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, username string, role types.UserRole) (string, error) {
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

	return token.SignedString([]byte(JWTConfig.SecretStr))
}

func ReadCtxEnv(c *gin.Context, targetStruct any) *types.CtxEnvUser {
	//解析的jwt信息：
	userIDAny, ok := c.Get("ID")
	if !ok {
		slog.Error("error when read ID from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read ID from env",
		})
		return nil
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		slog.Error("error when read ID from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read ID from env",
		})
		return nil
	}
	roleAny, ok := c.Get("Role")
	if !ok {
		slog.Error("error when read Role from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read Role from env",
		})
		return nil
	}
	role, ok := roleAny.(types.UserRole)
	if !ok {
		slog.Error("error when read Role from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read Role from env",
		})
		return nil
	}
	userNameAny, ok := c.Get("Username")
	if !ok {
		slog.Error("error when read Username from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read Username from env",
		})
		return nil
	}
	userName, ok := userNameAny.(string)
	if !ok {
		slog.Error("error when read Username from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read Username from env",
		})
		return nil
	}

	if err := c.ShouldBindJSON(targetStruct); err != nil {
		switch e := err.(type) {
		case validator.ValidationErrors:
			slog.Info("New lend request has no required value", "error", e)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "validation failed",
			})
			return nil
		case *json.UnmarshalTypeError:
			slog.Info("New lend request has dismatch type", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "value type dismatch",
			})
			return nil
		case *json.SyntaxError:
			slog.Info("New lend request has syntax error", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "json syntax error",
			})
			return nil
		default:
			slog.Error("error when unmarshal new lend json", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Internal server error",
			})
			return nil
		}
	}

	ctxEnv := &types.CtxEnvUser{
		ID:       userID,
		UserName: userName,
		Role:     role,
	}

	return ctxEnv
}
