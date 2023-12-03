package controller

import (
	"errors"

	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type PlayerController struct{}

func (pc *PlayerController) CreatePlayer(player models.Player) (models.Player, error) {
	if player.Uuid == "" {
		return player, errors.New("Player no tiene un UUID para crearse")
	}

	if err := database.DB.Create(&player).Error; err != nil {
		return player, err
	}
	return player, nil
}

func (pc *PlayerController) UpdatePLayer(player models.Player) (models.Player, error) {

	if err := database.DB.Save(&player).Error; err != nil {
		return player, err
	}
	return player, nil
}

func (pc *PlayerController) GetMyPlayer(idAccount uint) (players []models.Player, err error) {
	if err := database.DB.Where("account_id = ?", idAccount).Find(&players).Error; err != nil {
		return players, err
	}
	return players, err
}

func (pc *PlayerController) GetPlayer(name string) (players models.Player, err error) {

	if err := database.DB.Where("name = ?", name).First(&players).Error; err != nil {
		return players, err
	}
	return players, nil
}
