package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Hans-Kerman/go-book-lending/backend/global"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"github.com/Hans-Kerman/go-book-lending/backend/pkg"
	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	newUserRequest := &types.NewUser{}

	if err := c.ShouldBindJSON(newUserRequest); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok { //类型断言：是字段不符合要求
			slog.Info("error when bind json", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid field value",
			})
			return
		}
		slog.Error("error when bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error when bind json",
		})
		return
	}

	HashedPwd, err := pkg.HashPassword(newUserRequest.Password)
	if err != nil {
		slog.Error("error when hash password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error when hash password",
		})
		return
	}

	newUser := &models.User{
		Username: newUserRequest.UserName,
		Role:     models.Reader,
		Password: HashedPwd,
	}
	err = global.Db.Create(newUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) { //键名重复
			slog.Info("name conflict", "error", err)
			c.JSON(http.StatusConflict, gin.H{
				"error": "name conflict",
			})
			return
		}
		slog.Error("error when create database record", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unknown error",
		})
		return
	}

	token, err := pkg.GenerateJWT(newUser.ID, newUser.Username, newUser.Role)
	if err != nil {
		slog.Error("error when generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error when generate token",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}
