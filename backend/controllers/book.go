package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Hans-Kerman/go-book-lending/backend/global"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

func GetBookByISBN(c *gin.Context) {
	bookISBN := c.Param("isbn")

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
