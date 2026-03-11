package types

type BorrowRequire struct {
	BorrowReader uint   `json:"borrow_reader" binding:"required"`
	BookID       string `json:"book_id" binding:"required"`
}
