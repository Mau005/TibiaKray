package models

import "gorm.io/gorm"

type Comments struct {
	gorm.Model
	Comment   string
	TodaysID  uint
	AccountID uint
	Account   Account `gorm:"foreignKey:AccountID"`
}
