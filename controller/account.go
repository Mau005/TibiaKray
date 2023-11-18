package controller

import (
	"strings"

	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type AccountController struct{}

func (ac *AccountController) IsUserName(user string) bool {
	var acc models.Account
	if err := database.DB.Where("name = ?", user).First(&acc).Error; err != nil {
		return false
	}
	if acc.ID == 0 {
		return false
	}

	return true
}
func (ac *AccountController) GetAccountUser(user string) (*models.Account, error) {
	var acc *models.Account
	if err := database.DB.Where("name = ?", user).First(&acc).Error; err != nil {
		return nil, err
	}

	return acc, nil
}

func (ac *AccountController) CreateAccount(email, username, password string) (*models.Account, error) {
	var account models.Account
	account.Email = strings.ToLower(email)
	account.Name = strings.ToLower(username)

	var api ApiController
	account.Password = api.GenerateCryptPassword(password)

	if err := database.DB.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (ac *AccountController) GetAccount(email string) (acc *models.Account, err error) {
	if err = database.DB.Where("email = ?", strings.ToLower(email)).First(&acc).Error; err != nil {
		return acc, err
	}
	return acc, nil
}

func (ac *AccountController) Login(email, password string) (string, error) {
	acc, err := ac.GetAccount(email)
	if err != nil {
		acc, err = ac.GetAccountUser(email)
		if err != nil {
			return "", err
		}
	}

	var api ApiController
	err = api.CompareCryptPassword(acc.Password, password)
	if err != nil {
		return "", err
	}

	token, err := api.GenerateToken(acc)
	if err != nil {
		return "", err
	}

	return token, nil
}
