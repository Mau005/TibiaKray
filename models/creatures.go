package models

import "gorm.io/gorm"

type Creatures struct {
	gorm.Model
	Name       string `gorm:"unique"`
	ImagePath  string
	Health     uint
	Race       string
	Experience uint
	MaxDamage  uint
	Armor      uint
	Haste      uint
	Charm      uint
	Mitigation float32
	Dificulty  uint
	Locations  string
	Loot       string

	//Elements
	Physical uint
	Earth    uint
	Fire     uint
	Death    uint
	Energy   uint
	Holy     uint
	Ice      uint
	Healing  uint
	//End Elements

	PushObject     bool
	SummonConvince bool
	Pushable       bool
	Paralyzable    bool
}
