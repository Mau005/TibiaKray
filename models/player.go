package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Name            string `gorm:"unique;not null"`
	FormerNames     string
	Title           string
	Sex             string
	Vocation        string
	Level           string
	Achievement     string
	Points          string
	World           string
	FormerWorld     string
	Residence       string
	Comment         string
	House           string
	Guild           string
	LastLogin       string
	AccountStatus   string
	CountTimeActive *time.Time
	Uuid            string `gorm:"not null"`
	Status          bool   `gorm:"not null;default:0"`
	AccountID       *uint
	Account         Account `gorm:"foreignKey:AccountID"`
}
