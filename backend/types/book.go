package types

type NewBookInfo struct {
	Title          string `json:"title" binding:"required" example:"The Go Programming Language"`
	Author         string `json:"author" example:"Alan A. A. Donovan"`
	ISBN           string `json:"isbn" binding:"required" example:"978-0134190440"`
	Price          int    `json:"price" example:"40"`
	CoverPicBase64 []byte `json:"cover_pic_base64"`
	Available      *int   `json:"available" example:"10"`
}
