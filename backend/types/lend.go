package types

import "time"

type BorrowRequire struct {
	BorrowReader uint   `json:"borrow_reader" binding:"required"`
	BookID       string `json:"book_id" binding:"required"`
}

type LendRecordResponse struct {
	ID           uint       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	ReturnTime   *time.Time `json:"return_time,omitempty"`
	BorrowReader uint       `json:"borrow_reader"`
	BookID       string     `json:"book_id"`
}
