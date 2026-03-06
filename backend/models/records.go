package models

import "time"

type LendRecord struct {
	ID           uint       `gorm:"column:id;type:int;primaryKey;autoIncrement;uniqueIndex"`
	CreatedAt    time.Time  `gorm:"not null;index"`
	ReturnTime   *time.Time `gorm:"index"`
	BorrowReader uint       `gorm:"not null;index"`
	BookID       string     `gorm:"type:varchar(25);not null;index"` //使用ISBN书号
}
