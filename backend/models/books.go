package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title     string `gorm:"type:varchar(100);not null;index"`
	Author    string `gorm:"type:varchar(50);not null;index"`
	ISBN      string `gorm:"type:varchar(25);not null;uniqueIndex"`
	Price     int    `gorm:"default:0"` //存储以分为单位的整数作为(赔偿凭据)金额——备用
	Available int    //可借阅数量
	CoverURL  string `gorm:"type:varchar(255);"` // 存储图片的访问路径
}

func (book1 *Book) Equals(book2 *Book) bool {
	if book1.Title != book2.Title {
		return false
	}
	if book1.Author != book2.Author {
		return false
	}
	if book1.ISBN != book2.ISBN {
		return false
	}
	if book1.Price != book2.Price {
		return false
	}
	if book1.CoverURL != book2.CoverURL {
		return false
	}
	return true
}
