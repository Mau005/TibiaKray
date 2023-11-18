package models

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Title        string
	Content      string
	AccountID    *uint
	StatusNews   bool           //Check generate
	FileEncrypts []FileEncrypts `gorm:"foreignKey:NewsID"`
	Account      Account        `gorm:"foreignKey:AccountID"`
}

type NewsTicket struct {
	gorm.Model
	Title      string
	Content    string
	AccountID  *uint
	StatusNews bool
	Account    Account `gorm:"foreignKey:AccountID"`
}
