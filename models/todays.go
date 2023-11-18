package models

import "gorm.io/gorm"

type Todays struct {
	gorm.Model
	Title        string
	Description  string
	AccountID    uint
	Status       bool           `gorm:"default:false"`
	Account      Account        `gorm:"foreignKey:AccountID"`
	Comments     []Comments     `gorm:"foreignKey:TodaysID"`
	FileEncrypts []FileEncrypts `gorm:"foreignKey:TodaysID"`
}
