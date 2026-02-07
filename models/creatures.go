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
	Abilities  string
	Armor      uint
	Haste      uint
	Charm      uint
	Mitigation float32
	Dificulty  uint
	Locations  string
	Loot       string

	//Elements
	Physical int
	Earth    int
	Fire     int
	Death    int
	Energy   int
	Holy     int
	Ice      int
	Healing  int
	//End Elements

	Attack string

	PushObject     bool
	SummonConvince bool
	Pushable       bool
	Paralyzable    bool
}
