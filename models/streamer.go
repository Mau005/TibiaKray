package models

import "gorm.io/gorm"

type Streamers struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Title    string
	URL      string
	TypeUrl  string
	Lenguaje string
}
