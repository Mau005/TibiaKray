package controller

import (
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

func (vc *VotedController) GetVoted(idVoted uint, conditional string) (models.Voted, error) {
	var voted models.Voted
	if err := database.DB.Where(conditional, idVoted).First(&voted).Error; err != nil {
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

func (vc *VotedController) TodaysVoted(idTodays uint) models.Voted {
	voted, err := vc.GetVoted(idTodays, "todays_id")
	if err != nil {
		voted.TodaysID = &idTodays
		voted.Status = !voted.Status
		voted, _ = vc.CreatedVoted(voted)
		return voted
	}
	voted.Status = !voted.Status
	voted, _ = vc.UpdateVoted(voted)
	return voted
}
