package models

import "gorm.io/gorm"

type Todays struct {
	gorm.Model
	Title       string
	Description string
	AccountID   uint
	Views       uint       `gorm:"not null;default:0"`
	Status      bool       `gorm:"default:false"`
	Account     Account    `gorm:"foreignKey:AccountID"`
	Comments    []Comments `gorm:"foreignKey:TodaysID"`
	Files       []Files    `gorm:"foreignKey:TodaysID"`
	Voted       []Voted    `gorm:"foreignKey:TodaysID"`
}
