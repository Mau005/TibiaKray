package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name          string
	FormerNames   string
	Title         string
	Sex           string
	Vocation      string
	Level         string
	Achievement   string
	Points        string
	World         string
	FormerWorld   string
	Residence     string
	Comment       string
	House         string
	Guild         string
	LastLogin     string
	AccountStatus string
}
