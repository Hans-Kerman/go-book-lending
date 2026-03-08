package pkg

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindNewUser(c *gin.Context, newUserRequest *types.NewUser) bool {
	if err := c.ShouldBindJSON(newUserRequest); err != nil {
		switch e := err.(type) {
		case validator.ValidationErrors:
			slog.Info("validation failed when bind json", "errors", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid filed value",
			})
		case *json.SyntaxError:
			slog.Info("json syntax error when bind json", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Malformed JSON",
			})
		case *json.UnmarshalTypeError:
			slog.Info("json type dismatch", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Field type mismatch",
			})
		default:
			slog.Error("error when bind user json", "error", e)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return false
	}
	return true
}
