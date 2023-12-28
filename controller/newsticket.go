package controller

import (
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type NewsTicketController struct{}

func (nc *NewsTicketController) GetNewsTicket() (new []models.NewsTicket, err error) {
	if err = database.DB.Order("created_at desc").Where("status_news = 1").Limit(5).Find(&new).Error; err != nil {
		return
	}
	return
}
func (nc *NewsTicketController) GetIDNewsTicket(idTicket uint) (new models.NewsTicket, err error) {
	if err = database.DB.Where("id = ?", idTicket).First(&new).Error; err != nil {
		return
	}
	return

}
func (nc *NewsTicketController) GetAllNewsTicket() (new []models.NewsTicket, err error) {
	if err = database.DB.Order("created_at desc").Find(&new).Error; err != nil {
		return
	}
	return

}

func (nc *NewsTicketController) SaveNewsTicket(new models.NewsTicket) (models.NewsTicket, error) {
	if err := database.DB.Save(&new).Error; err != nil {
		return new, err
	}
	return new, nil
}

func (nc *NewsTicketController) CreateNewsTicket(new models.NewsTicket) (models.NewsTicket, error) {

	if err := database.DB.Create(&new).Error; err != nil {
		return new, err
	}
	return new, nil
}
