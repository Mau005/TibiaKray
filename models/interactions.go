package models

import "gorm.io/gorm"

type Interactions struct {
	gorm.Model
	TodaysID              *uint
	NewsID                *uint
	ForumThreadsContentID *uint
	AccountID             uint
}
