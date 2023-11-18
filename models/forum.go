package models

import "gorm.io/gorm"

type AbstracForum struct {
	Name      string
	Content   string
	IconPath  string
	AccountID uint
	Account   Account `gorm:"foreignkey:AccountID"`
}

type ForumThreadsContent struct {
	gorm.Model
	AbstracForum
	ForumThreadsID uint
	AccountID      uint
	Account        Account `gorm:"foreignkey:AccountID"`
}

type ForumThreads struct {
	gorm.Model
	AbstracForum
	ForumID             uint
	ForumThreadsContent []ForumThreadsContent `gorm:"foreignKey:ForumThreadsID"`
}

type Forum struct {
	gorm.Model
	AbstracForum
	ForumThreads []ForumThreads `gorm:"foreignKey:ForumID"`
}
