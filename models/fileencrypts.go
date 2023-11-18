package models

import "gorm.io/gorm"

type FileEncrypts struct {
	gorm.Model
	PathEncrypt string `gorm:"unique" json:"pathencrypt"`
	PathOrigin  string `gorm:"not null" json:"-"`
	NewsID      *uint
	TodaysID    *uint
}
