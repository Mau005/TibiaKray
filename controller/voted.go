package controller

import (
	"log"

	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type VotedController struct{}

func (vc *VotedController) CreatedVoted(voted models.Voted) (models.Voted, error) {
	if err := database.DB.Create(&voted).Error; err != nil {
		return voted, err
	}
	return voted, nil

}

func (vc *VotedController) GetVotedAccountTodays(idTodays, accountID uint) (models.Voted, error) {
	var voted models.Voted
	if err := database.DB.Where("todays_id = ? AND account_id = ?", idTodays, accountID).First(&voted).Error; err != nil {
		return voted, err
	}
	return voted, nil

}

func (vc *VotedController) UpdateVoted(voted models.Voted) (models.Voted, error) {
	if err := database.DB.Save(&voted).Error; err != nil {
		return voted, err
	}
	return voted, nil
}

func (vc *VotedController) TodaysVoted(idTodays uint, sm models.StructModel) (models.Voted, error) {
	var voted models.Voted
	var accManager AccountController
	account, err := accManager.GetAccount(sm.Email)
	if err != nil {
		log.Println("Error de usuario: ", err)
		return voted, err
	}

	if err := database.DB.Where("todays_id = ? AND account_id = ?", idTodays, account.ID).First(&voted).Error; err != nil {
		voted.AccountID = account.ID
		voted.Status = !voted.Status
		voted.TodaysID = &idTodays
		voted, err = vc.CreatedVoted(voted)
		if err != nil {
			log.Println("Error al Crear: ", err)
			return voted, err
		}
		return voted, nil
	}

	voted.Status = !voted.Status
	voted, err = vc.UpdateVoted(voted)
	if err != nil {
		log.Println("Error al actualizar voted: ", err)
		return voted, err
	}
	return voted, nil
}
