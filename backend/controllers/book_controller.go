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
	"gorm.io/gorm/clause"
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

// @Summary Get books by page
// @Description Get a paginated list of books
// @Tags Books
// @Produce  json
// @Param   page      query   int     false  "Page number"
// @Param   page_size query   int     false  "Page size"
// @Success 200 {object} map[string]interface{} "{"total": 1, "page": 1, "page_size": 20, "totalPages": 1, "books": []models.Book}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Router /public/books [get]
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

// @Summary Get a book by ISBN
// @Description Get a single book by its ISBN
// @Tags Books
// @Produce  json
// @Param   isbn     path    string     true        "Book ISBN"
// @Success 200 {object} map[string]models.Book "{"book": models.Book}"
// @Failure 404 {object} map[string]string "{"error": "Book not found"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Router /public/book/{isbn} [get]
func GetBookByURL(c *gin.Context) {
	bookISBN := c.Param("isbn")
	getBookByISBN(c, bookISBN)
}

// @Summary Add a new book
// @Description Add a new book to the database (Admin only)
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param   book     body    types.NewBookInfo     true        "New Book Info"
// @Success 200 {object} map[string]models.Book "{"data": models.Book}"
// @Failure 400 {object} map[string]string "{"error": "error_message"}"
// @Failure 422 {object} map[string]string "{"error": "validation failed"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /admin/book [post]
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

// @Summary Update a book
// @Description Update an existing book's information (Admin only)
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param   book     body    types.NewBookInfo     true        "Book Info to Update"
// @Success 200 {object} map[string]models.Book "{"data": models.Book}"
// @Failure 400 {object} map[string]string "{"error": "error_message"}"
// @Failure 404 {object} map[string]string "{"error": "Assigned book not found"}"
// @Failure 422 {object} map[string]string "{"error": "validation failed"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /admin/book [put]
func UpdateBook(c *gin.Context) {
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

	storedBook := &models.Book{}

	err := global.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("isbn = ?", newBookInfo.ISBN).Clauses(clause.Locking{Strength: "UPDATE"}).First(storedBook)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				slog.Info("client try to update not-exist book", "isbn", newBookInfo.ISBN)
				return result.Error
			}
			slog.Error("Unknown error when query book data", "error", result.Error)
			return result.Error
		}

		if storedBook.EqualsDTO(newBookInfo) {
			slog.Info("client update but no diffs", "isbn", newBookInfo.ISBN)
			return nil
		}

		storedBook.Title = newBookInfo.Title
		storedBook.Author = newBookInfo.Author
		storedBook.Price = newBookInfo.Price
		if newBookInfo.CoverPicBase64 != nil {
			storedBook.CoverURL = pkg.ParsePicName(newBookInfo.CoverPicBase64)
		}

		if newBookInfo.Available != nil {
			storedBook.Available = *newBookInfo.Available
		}
		result = tx.Save(storedBook)

		if result.Error != nil {
			slog.Error("error when update book", "error", result.Error)
			return result.Error
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Assigned book not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": storedBook,
	})
}

// @Summary Delete a book
// @Description Delete a book by its ISBN (Admin only)
// @Tags Admin
// @Produce  json
// @Param   isbn     path    string     true        "Book ISBN"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "{"error": "Target book not exist"}"
// @Failure 500 {object} map[string]string "{"error": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /admin/book/del/{isbn} [delete]
func DelBook(c *gin.Context) {
	isbn := c.Param("isbn")

	result := global.Db.Where("isbn = ?", isbn).Delete(&models.Book{})

	if result.Error != nil {
		slog.Error("Unkonwn error when delete book", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		slog.Info("client try to delete no-exist book", "isbn", isbn)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Target book not exist",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
