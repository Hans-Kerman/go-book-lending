package middlewares

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Hans-Kerman/go-book-lending/backend/config"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"github.com/Hans-Kerman/go-book-lending/backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ParseJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			slog.Info("client request unauthorized")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized"},
			)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			slog.Info("client request with authorization syntax error")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization syntax error",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		parsedClaims := &pkg.UserClaims{}
		_, err := jwt.ParseWithClaims(tokenString, parsedClaims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.AppConfig.JWT.SecretStr), nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				// token 格式错误
				slog.Info("client auth with malformed jwt", "error", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "jwt malformed",
				})
				ctx.Abort()
				return
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				// 签名无效
				slog.Info("client auth with invalid signature", "error", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "jwt signature invalid",
				})
				ctx.Abort()
				return
			case errors.Is(err, jwt.ErrTokenExpired):
				// token 过期，可以返回 401 并提示重新登录
				slog.Info("client auth with expired jwt", "error", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "jwt expired",
				})
				ctx.Abort()
				return
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				// token 还未生效
				slog.Info("client auth with not valid yet jwt", "error", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "jwt not valid yet",
				})
				ctx.Abort()
				return
			default:
				slog.Error("unknown error when auth", "error", err)
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unknown error",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Set("ID", parsedClaims.ID)
		ctx.Set("Username", parsedClaims.UserName)
		ctx.Set("Role", parsedClaims.Role)

		ctx.Next()
	}
}

func CheckAdminRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := ctx.Get("Role")
		userID, _ := ctx.Get("ID")
		if role == models.Admin {
			ctx.Next()
		} else {
			slog.Info("reader client require admin api", "user_id", userID, "role", role, "api", ctx.Request.URL)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
		}
	}
}
