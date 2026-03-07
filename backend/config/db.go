package config

import (
	"fmt"

	"github.com/Hans-Kerman/go-book-lending/backend/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDataBase() error {
	dbConfig := AppConfig.Pgsql
	dsn := fmt.Sprintf("host=%s user=%s password=%s "+
		"dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		&dbConfig.Host, &dbConfig.User, &dbConfig.Password,
		&dbConfig.DbName, dbConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error when init database: %w", err)
	}

	global.Db = db
	return nil
}
