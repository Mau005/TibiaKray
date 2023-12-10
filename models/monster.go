package models

import "gorm.io/gorm"

type Monster struct {
	gorm.Model
	Name       string `gorm:"unique, not null"`
	ImagePath  string
	Experience uint
	Abilities  string
	Pushable   bool
	MaxDamage  uint
}

type Elements struct {
	gorm.Model
	Name  string
	Count int
}
