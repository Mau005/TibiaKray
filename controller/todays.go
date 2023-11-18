package controller

import (
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type TodaysController struct{}

func (tc *TodaysController) GetTodays(id uint) (today models.Todays, err error) {

	if err := database.DB.Where("id = ?", id).First(&today).Error; err != nil {
		return today, err
	}

	return today, err
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
