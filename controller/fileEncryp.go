package controller

import (
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type FileEncryptsController struct{}

func (fec *FileEncryptsController) CreateEncrypFile(fileEncr models.FileEncrypts) (models.FileEncrypts, error) {

	if err := database.DB.Create(&fileEncr).Error; err != nil {
		return fileEncr, err
	}

	return fileEncr, nil

}

func (fec *FileEncryptsController) GetEncrypFile(hash string) (string, error) {
	var enc models.FileEncrypts
	if err := database.DB.Where("path_encrypt = ?", hash).First(&enc).Error; err != nil {
		return enc.PathOrigin, err
	}
	return enc.PathOrigin, nil
}
