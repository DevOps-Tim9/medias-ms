package entity

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	Url string `gorm:"not null;default:null"`
	Tbl string `gorm:"-"`
}
