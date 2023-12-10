package models

import "gorm.io/gorm"

type Files struct {
	gorm.Model
	PathEncrypt string `gorm:"unique" json:"pathencrypt"`
	PathOrigin  string `gorm:"not null" json:"-"`
	PathConsume string
	NewsID      *uint
	TodaysID    *uint
	Todays      *Todays `gorm:"foreignKey:TodaysID"`
	TypeFile    string
}
