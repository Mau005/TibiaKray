package controller

import (
	"errors"

	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
)

type StreamerController struct{}

func (sc *StreamerController) CreateStremer(streamer models.Streamers) (models.Streamers, error) {
	if streamer.Name == "" {
		return streamer, errors.New("EL nombre no puede estar sin datos")
	}

	if err := database.DB.Create(&streamer).Error; err != nil {
		return streamer, err
	}
	return streamer, nil
}

func (sc *StreamerController) SaveStremer(streamer models.Streamers) (models.Streamers, error) {
	if err := database.DB.Save(&streamer).Error; err != nil {
		return streamer, err
	}
	return streamer, nil
}

func (sc *StreamerController) GetStreamers() (streamers []models.Streamers, err error) {
	if err = database.DB.Find(&streamers).Error; err != nil {
		return
	}
	return
}

func (sc *StreamerController) GetIdStreamers(idStreamers uint) (streamers []models.Streamers, err error) {
	if err = database.DB.Where("id = ?", idStreamers).First(&streamers).Error; err != nil {
		return
	}
	return
}

func (sc *StreamerController) DeleteStreamers(idStreamers uint) error {
	streamers, err := sc.GetIdStreamers(idStreamers)
	if err != nil {
		return err
	}
	if err = database.DB.Delete(&streamers).Error; err != nil {
		return err
	}
	return nil
}
