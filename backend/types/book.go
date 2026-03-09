package types

type NewBookInfo struct {
	Title          string `json:"title" binding:"required"`
	Author         string `json:"author"`
	ISBN           string `json:"isbn" binding:"required"`
	Price          int    `json:"price"`
	CoverPicBase64 []byte `json:"cover_pic_base64"`
}
