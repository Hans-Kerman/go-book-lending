package controllers

import (
	"errors"
	"log"
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
		log.Println("error when bind json: " + err.Error())
		if _, ok := err.(validator.ValidationErrors); ok { //类型断言：是字段不符合要求
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid field value",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error when bind json",
		})
		return
	}

	HashedPwd, err := pkg.HashPassword(newUserRequest.Password)
	if err != nil {
		log.Println("error when hash password: " + err.Error())
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
			log.Println("name conflict: " + err.Error())
			c.JSON(http.StatusConflict, gin.H{
				"error": "name conflict",
			})
			return
		}
		log.Println("error when create database record: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unknown error",
		})
		return
	}

	token, err := pkg.GenerateJWT(newUser.ID, newUser.Username, newUser.Role)
	if err != nil {
		log.Println("error when generate token: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error when generate token",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}
