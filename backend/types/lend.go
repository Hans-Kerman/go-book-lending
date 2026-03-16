package types

import "time"

type BorrowRequire struct {
	BorrowReader uint   `json:"borrow_reader" binding:"required" example:"1"`
	BookID       string `json:"book_id" binding:"required" example:"978-0134190440"`
}

type LendRecordResponse struct {
	ID           uint       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	ReturnTime   *time.Time `json:"return_time,omitempty"`
	BorrowReader uint       `json:"borrow_reader"`
	BookID       string     `json:"book_id"`
}
