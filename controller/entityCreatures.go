package controller

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gocolly/colly/v2"
)

type EntitysCreatures struct{}

func (ec *EntitysCreatures) CollectorCreature() error {
	c := colly.NewCollector()
	var monsters []models.Monster
	c.OnHTML("div[style='width: 100px; height: 110px; margin: 0px; float: left;']", func(e *colly.HTMLElement) {
		// Encuentra el enlace dentro del div
		// Encuentra la etiqueta img dentro del div
		imgSrc := e.ChildAttr("img", "src")
		// Encuentra el texto dentro del div
		name := strings.TrimSpace(e.ChildText("div"))

		monsters = append(monsters, models.Monster{Name: name, ImagePath: imgSrc})
	})

	err := c.Visit(configuration.TIBIA_MONSTER)
	if err != nil {
		log.Println(err)
	}
	for _, value := range monsters {
		var api ApiController
		sliceNames := strings.Split(value.ImagePath, "/")
		nameImage := sliceNames[len(sliceNames)-1]
		pathSave := fmt.Sprintf(configuration.PATH_STATIC_CREATURES, nameImage)
		api.DownloadImage(value.ImagePath, pathSave)
		value.ImagePath = pathSave
		_, err := ec.CreateMonster(value)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (ec *EntitysCreatures) CollectorBosses() error {
	c := colly.NewCollector()
	var bosses []models.Bosses
	c.OnHTML("div[style='width: 100px; height: 110px; margin: 0px; float: left;']", func(e *colly.HTMLElement) {
		// Encuentra el enlace dentro del div
		// Encuentra la etiqueta img dentro del div
		imgSrc := e.ChildAttr("img", "src")
		// Encuentra el texto dentro del div
		name := strings.TrimSpace(e.ChildText("div"))

		bosses = append(bosses, models.Bosses{Monster: models.Monster{Name: name, ImagePath: imgSrc}})
	})

	err := c.Visit(configuration.TIBIA_BOSSES)
	if err != nil {
		log.Println(err)
	}
	for _, value := range bosses {
		var api ApiController
		sliceNames := strings.Split(value.ImagePath, "/")
		nameImage := sliceNames[len(sliceNames)-1]
		pathSave := fmt.Sprintf(configuration.PATH_STATIC_BOSSES, nameImage)
		api.DownloadImage(value.ImagePath, pathSave)
		value.ImagePath = pathSave
		_, err := ec.CreateBosses(value)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (ec *EntitysCreatures) CreateMonster(monster models.Monster) (models.Monster, error) {
	if err := database.DB.Create(&monster).Error; err != nil {
		return monster, err
	}
	return monster, nil
}

func (ec *EntitysCreatures) CreateBosses(boss models.Bosses) (models.Bosses, error) {
	if err := database.DB.Create(&boss).Error; err != nil {
		return boss, err
	}
	return boss, nil
}
