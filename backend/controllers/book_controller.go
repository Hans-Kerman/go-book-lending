package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Hans-Kerman/go-book-lending/backend/global"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"github.com/Hans-Kerman/go-book-lending/backend/pkg"
	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func getBookByISBN(c *gin.Context, bookISBN string) {
	expectedBook := &models.Book{
		ISBN: bookISBN,
	}

	if err := global.Db.Model(&models.Book{}).Where(expectedBook).First(expectedBook).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("client request unstored book")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Book not found",
			})
			return
		}
		slog.Error("error when query book in database", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"book": expectedBook,
	})
}

func GetBooksByPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	//构建查询
	query := global.Db.Model(&models.Book{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		slog.Error("error when count querys", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	var books []models.Book
	err := query.Order("id asc").Offset(offset).Limit(pageSize).Find(&books).Error
	if err != nil {
		slog.Error("error when query data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		"books":      books,
	})
}

func GetBookByURL(c *gin.Context) {
	bookISBN := c.Param("isbn")
	getBookByISBN(c, bookISBN)
}

func PostNewBook(c *gin.Context) {
	newBookInfo := &types.NewBookInfo{}
	if err := c.ShouldBindJSON(newBookInfo); err != nil {
		switch e := err.(type) {
		case validator.ValidationErrors:
			slog.Info("New book request has no required value", "error", e)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "validation failed",
			})
			return
		case *json.UnmarshalTypeError:
			slog.Info("New book request has dismatch type", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "value type dismatch",
			})
			return
		case *json.SyntaxError:
			slog.Info("New book request has syntax error", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "json syntax error",
			})
			return
		default:
			slog.Error("error when unmarshal new book json", "error", e)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Internal server error",
			})
			return
		}
	}

	newBook := &models.Book{
		Title:     newBookInfo.Title,
		Author:    newBookInfo.Author,
		ISBN:      newBookInfo.ISBN,
		Price:     newBookInfo.Price,
		Available: 0,
		CoverURL:  pkg.ParsePicName(newBookInfo.CoverPicBase64),
	}

	if err := global.Db.Model(&models.Book{}).Create(newBook).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			slog.Info("New book to create already exist", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Book already exists",
			})
			return
		}
		slog.Error("Internal server error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": newBook,
	})
}
