package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name      string `gorm:"unique" json:"name"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"-"`
	EndPremmy *time.Time
	Groups    uint8  `gorm:"default:0" json:"groups"`
	Languaje  string `gorm:"default:en"`
}
