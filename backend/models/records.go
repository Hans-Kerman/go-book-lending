package models

import (
	"time"

	"github.com/Hans-Kerman/go-book-lending/backend/types"
)

type LendRecord struct {
	ID           uint       `gorm:"column:id;type:int;primaryKey;autoIncrement;uniqueIndex"`
	CreatedAt    time.Time  `gorm:"not null;index"`
	ReturnTime   *time.Time `gorm:"index"`
	BorrowReader uint       `gorm:"not null;index"`
	BookID       string     `gorm:"type:varchar(25);not null;index"` //使用ISBN书号
}

func (record *LendRecord) ConvertResp() types.LendRecordResponse {
	return types.LendRecordResponse{
		ID:           record.ID,
		CreatedAt:    record.CreatedAt,
		ReturnTime:   record.ReturnTime,
		BorrowReader: record.BorrowReader,
		BookID:       record.BookID,
	}
}
