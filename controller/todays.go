package controller

import (
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"gorm.io/gorm"
)

type TodaysController struct{}

func (tc *TodaysController) GetToday(id uint) (today models.Todays, err error) {
	tx := database.DB.Begin()

	if tx.Error != nil {
		return today, tx.Error
	}

	if err := tx.Preload("Files").Preload("Voted", "status = ?", 1).Preload("Comments.Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, access")
	}).Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, access")
	}).Where("id = ?", id).First(&today).Error; err != nil {

		tx.Rollback()
		return today, err
	}

	if err := tx.Model(&today).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
		tx.Rollback()
		return today, err
	}

	// Commit de la transacción
	tx.Commit()

	return today, nil
}

func (tc *TodaysController) GetTodayView() (todays []models.Todays, err error) {

	if err = database.DB.Preload("Files").Where("status = 1").Order("created_at desc").Find(&todays).Error; err != nil {
		return todays, err
	}
	if err := database.DB.Model(&todays).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
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

	if len(todays) > 0 {
		// Suponemos que el primer elemento de la respuesta contiene los datos deseados
		firstPhoto := todays[0]

		// Incrementar el contador
		if err := database.DB.Model(&firstPhoto).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
			return todays, err
		}
	}

	return todays, nil
}

func (tc *TodaysController) GetTodayPage(page int) (photos []models.Todays, err error) {
	offset := (page - 1) * 4
	limit := 4

	if err := database.DB.Preload("Files").Preload("Comments.Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, access")
	}).Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email, access")
	}).Order("created_at desc").Where("status = 1").Offset(offset).Limit(limit).Find(&photos).Error; err != nil {
		return photos, err
	}

	if len(photos) > 0 {
		// Suponemos que el primer elemento de la respuesta contiene los datos deseados
		firstPhoto := photos[0]

		// Incrementar el contador
		if err := database.DB.Model(&firstPhoto).UpdateColumn("Views", gorm.Expr("views + ?", 1)).Error; err != nil {
			return photos, err
		}
	}

	return photos, err
}
