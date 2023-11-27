package models

import (
	"time"
)

type SharedLoot struct {
	DateStart time.Time
	DateEnd   time.Time
	Time      time.Time
	LootType  string
	Loot      int
	Supplies  int
	Balance   int
	Leader    CharacterShared
	Character []CharacterShared
}

type CharacterShared struct {
	Name     string
	Loot     int
	Supplies int
	Balance  int
	Damage   int
	Healing  int
}
