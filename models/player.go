package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	AccountId     uint
	Name          string
	Title         string
	Sex           string
	Vocation      string
	Level         string
	Achievement   string
	Points        string
	World         string
	Residence     string
	Guild         string
	Login         string
	CETComment    string
	AccountStatus string
}
