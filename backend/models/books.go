package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title     string `gorm:"type:varchar(100);not null;index"`
	Author    string `gorm:"type:varchar(50);not null;index"`
	ISBN      string `gorm:"type:varchar(25);not null;uniqueIndex"`
	Price     int    `gorm:"default:0"` //存储以分为单位的整数作为(赔偿凭据)金额——备用
	Available int    //可借阅数量
}
