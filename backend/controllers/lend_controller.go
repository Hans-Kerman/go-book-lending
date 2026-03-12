package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Hans-Kerman/go-book-lending/backend/global"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"github.com/Hans-Kerman/go-book-lending/backend/pkg"
	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var bookNumInsufficientErr = errors.New("Insufficient book stock quantity")
var bookDataDisappearErr = errors.New("Unknown error cause book data disappeared")

func LendBook(c *gin.Context) {
	newLendRequire := &types.BorrowRequire{}
	env := pkg.ReadCtxEnv(c, newLendRequire)
	if env == nil {
		//ReadCtxEnv()中已经写入错误响应
		return
	}

	if env.Role != types.Admin && env.ID != newLendRequire.BorrowReader {
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
	newReturnRequire := &types.BorrowRequire{}
	env := pkg.ReadCtxEnv(c, newReturnRequire)
	if env == nil {
		//ReadCtxEnv()中已经写入错误响应
		return
	}

	if env.Role != types.Admin &&
		env.ID != newReturnRequire.BorrowReader {
		slog.Info("client try to return book with ID mismatch",
			"sender", env.ID,
			"req_user", newReturnRequire.BorrowReader,
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only admin can return for others",
		})
		return
	}

	targetRecord := &models.LendRecord{
		BorrowReader: newReturnRequire.BorrowReader,
		BookID:       newReturnRequire.BookID,
		ReturnTime:   nil,
	}
	//构造事务
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		//查询待归还的记录
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("borrow_reader = ? AND book_id = ? AND return_time IS NULL",
				targetRecord.BorrowReader, targetRecord.BookID).
			First(targetRecord).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Info("Target book record does not exist")
				return err
			}
			slog.Error("error when query record for returning", "error", err)
			return err
		}

		//更新图书存量
		result := tx.Model(&models.Book{}).Where("isbn = ?", targetRecord.BookID).
			Update("available", gorm.Expr("available + 1"))
		if result.Error != nil {
			slog.Error("error when update book num", "error", err)
			return err
		}
		if result.RowsAffected == 0 {
			slog.Error("error when get book data", "error", err, "book_isbn", targetRecord.BookID)
			return bookDataDisappearErr
		}

		//更新归还时间
		err = tx.Model(&models.LendRecord{}).
			Where("id = ?", targetRecord.ID).
			Update("return_time", time.Now()).Error
		if err != nil {
			slog.Error("error when update return time", "error", err)
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Borrow record not found",
			})
			return
		}
		if errors.Is(err, bookDataDisappearErr) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal serverl error",
		})
		return
	}

	now := time.Now()
	targetRecord.ReturnTime = &now
	c.JSON(http.StatusOK, gin.H{
		"data": targetRecord,
	})
}

func GetUserRecord(c *gin.Context) {
	userIDAny, ok := c.Get("ID")
	if !ok {
		slog.Error("error when read ID from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read ID from env",
		})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		slog.Error("error when read ID from parsed jwt")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "can't read ID from env",
		})
		return
	}

	records := []models.LendRecord{}
	res := global.Db.Model(&models.LendRecord{}).Where("borrow_reader = ?", userID).Find(&records)
	if res.Error != nil {
		slog.Error("error when query records for user", "error", res.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	recordsResp := []types.LendRecordResponse{}
	for i := range records {
		recordsResp = append(recordsResp, records[i].ConvertResp())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": recordsResp,
	})
}
