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
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	newUserRequest := &types.NewUser{}
	if ok := pkg.BindNewUser(c, newUserRequest); !ok {
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
		Role:     types.Reader,
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

func Login(c *gin.Context) {
	newUserRequest := &types.NewUser{}
	if ok := pkg.BindNewUser(c, newUserRequest); !ok {
		return
	}

	expectedUser := &models.User{
		Username: newUserRequest.UserName,
	}
	if result := global.Db.Take(expectedUser); result.Error != nil {
		err := result.Error
		switch err {
		case gorm.ErrRecordNotFound:
			slog.Info("user not found in database", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "username or password incorrect",
			})
			return
		default:
			slog.Error("error when select database", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}
	}

	if err := pkg.CheckPassword(newUserRequest.Password, expectedUser.Password); err != nil {
		slog.Info("password validate failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "username or password incorrect",
		})
		return
	}

	token, err := pkg.GenerateJWT(expectedUser.ID, expectedUser.Username, expectedUser.Role)
	if err != nil {
		slog.Error("error when generate jwt", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	slog.Info("user login succeeded")
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
