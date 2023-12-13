package models

import "gorm.io/gorm"

type Creatures struct {
	gorm.Model
	Name       string `gorm:"unique"`
	ImagePath  string
	Experience uint
	Abilities  string
	Pushable   bool
	MaxDamage  uint
}
