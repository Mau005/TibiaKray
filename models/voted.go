package models

import "gorm.io/gorm"

type Voted struct {
	gorm.Model
	Status    bool
	TodaysID  *uint
	CommentID *uint
	AccountID uint
	Todays    Todays   `gorm:"foreignKey:TodaysID"`
	Comments  Comments `gorm:"foreignKey:CommentID"`
	Account   Account  `gorm:"foreignKey:AccountID"`
}
