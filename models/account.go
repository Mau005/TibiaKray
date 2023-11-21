package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name         string `gorm:"unique" json:"name"`
	Email        string `gorm:"unique" json:"email"`
	Password     string `json:"-"`
	EndPremmy    *time.Time
	Access       uint8          `gorm:"default:0" json:"access"`
	Languaje     string         `gorm:"default:en"`
	StreamMode   bool           `gorm:"default:false"`
	Interactions []Interactions `gorm:"foreignKey:AccountID"`
}
