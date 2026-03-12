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
	"gorm.io/gorm/clause"
)

var bookNumInsufficientErr = errors.New("Insufficient book stock quantity")

func LendBook(c *gin.Context) {
	newLendRequire := &types.BorrowRequire{}
	env := pkg.ReadCtxEnv(c, newLendRequire)
	if env == nil {
		//ReadCtxEnv()中已经写入错误响应
		return
	}

	if env.Role != models.Admin && env.ID != newLendRequire.BorrowReader {
		slog.Info("client try borrow for other without admin", "user_id", env.ID)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only admin can borrow for others",
		})
		return
	}

	dbLendRecord := &models.LendRecord{}

	//开始数据库事务
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		//构造查询
		targetBook := &models.Book{
			ISBN: newLendRequire.BookID,
		}

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("isbn = ?", targetBook.ISBN).
			First(targetBook).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Info("client require book not found", "user_id", env.ID, "book_isbn", targetBook.ISBN)
				return err
			}
			slog.Error("unknown error when query required book", "error", err)
			return err
		}

		//数量不足返回
		if targetBook.Available <= 0 {
			slog.Info("client try to borrow a 0available book", "user_id", env.ID, "book_name", targetBook.Title)
			slog.Warn("Book insufficient", "book_name", targetBook.Title)
			return bookNumInsufficientErr
		}

		//构造写入借书记录
		newLendRecord := &models.LendRecord{
			ReturnTime:   nil,
			BorrowReader: newLendRequire.BorrowReader,
			BookID:       newLendRequire.BookID,
		}

		if err := tx.Create(newLendRecord).Error; err != nil {
			slog.Error("error when create database record", "error", err)
			return err
		}

		//减少库存记录
		if err := tx.Model(targetBook).Where("available >= ?", 1).
			Update("available", gorm.Expr("available - ?", 1)).Error; err != nil {
			return err
		}

		dbLendRecord = newLendRecord
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Required book not found",
			})
			return
		}
		if errors.Is(err, bookNumInsufficientErr) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Insufficient book stock quantity",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": dbLendRecord,
	})
}

func ReturnBook(c *gin.Context) {
}
