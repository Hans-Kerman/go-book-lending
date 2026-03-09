package config

import (
	"fmt"

	"github.com/Hans-Kerman/go-book-lending/backend/global"
	"github.com/Hans-Kerman/go-book-lending/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDataBase() error {
	dbConfig := AppConfig.Pgsql
	dsn := fmt.Sprintf("host=%s user=%s password=%s "+
		"dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.DbName, dbConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return fmt.Errorf("error when init database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("error when AutoMigrate database: %w", err)
	}
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		return fmt.Errorf("error when AutoMigrate database: %w", err)
	}
	if err := db.AutoMigrate(&models.LendRecord{}); err != nil {
		return fmt.Errorf("error when AutoMigrate database: %w", err)
	}

	global.Db = db
	return nil
}
