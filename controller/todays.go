package controller

import (
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"gorm.io/gorm"
)

type TodaysController struct{}

func (tc *TodaysController) GetToday(id uint) (today models.Todays, err error) {
	if err := database.DB.Preload("Files").Preload("Comments.Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, access")
	}).Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, access")
	}).Where("id = ?", id).First(&today).Error; err != nil {
		return today, err
	}

	return today, err
}

func (tc *TodaysController) GetTodayView() (todays []models.Todays, err error) {

	if err = database.DB.Preload("Files").Where("status = 1").Order("created_at desc").Find(&todays).Error; err != nil {
		return todays, err
	}
	return todays, err
}

func (tc *TodaysController) SaveTodays(today models.Todays) (models.Todays, error) {

	if err := database.DB.Save(&today).Error; err != nil {
		return today, err
	}

	return today, nil
}

func (tc *TodaysController) CreateTodays(today models.Todays) (models.Todays, error) {

	if err := database.DB.Create(&today).Error; err != nil {
		return today, err
	}

	return today, nil
}

func (tc *TodaysController) GetAllTodaysStatus(status int) (todays []models.Todays, err error) {
	if err = database.DB.Preload("Account").Where("status = ? ", status).Find(&todays).Error; err != nil {
		return todays, err
	}
	return todays, err
}

func (tc *TodaysController) GetTodaysLobby() (todays []models.Todays, err error) {

	if err := database.DB.Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("Files").Order("created_at  desc").Limit(3).Where("status = 1").Find(&todays).Error; err != nil {
		return todays, err
	}
	return todays, nil
}
