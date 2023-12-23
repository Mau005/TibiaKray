package controller

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gocolly/colly/v2"
)

type EntitysCreatures struct{}

func (ec *EntitysCreatures) CollectorCreature() error {
	c := colly.NewCollector()
	var monsters []models.Creatures
	c.OnHTML("div[style='width: 100px; height: 110px; margin: 0px; float: left;']", func(e *colly.HTMLElement) {
		// Encuentra el enlace dentro del div
		// Encuentra la etiqueta img dentro del div
		imgSrc := e.ChildAttr("img", "src")
		// Encuentra el texto dentro del div
		name := strings.TrimSpace(e.ChildText("div"))

		monsters = append(monsters, models.Creatures{Name: name, ImagePath: imgSrc})
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

		bosses = append(bosses, models.Bosses{Creatures: models.Creatures{Name: strings.Trim(name, " "), ImagePath: imgSrc}})
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
		err := api.DownloadImage(value.ImagePath, pathSave)
		log.Println(err)
		value.ImagePath = pathSave
		_, err = ec.CreateBosses(value)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (ec *EntitysCreatures) GetCreatures() (creatures []models.Creatures, err error) {
	if err = database.DB.Find(&creatures).Error; err != nil {
		return
	}
	return
}
func (ec *EntitysCreatures) GetIdCreatures(idCreature uint) (creatures models.Creatures, err error) {
	if err = database.DB.Where("id = ?", idCreature).First(&creatures).Error; err != nil {
		return
	}
	return
}

func (ec *EntitysCreatures) GetThreeCreaturesRand() (creatures []models.Creatures, err error) {
	var allCreatures []models.Creatures
	database.DB.Find(&allCreatures)

	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	return
}
func (ec *EntitysCreatures) GetThreeBossRand() (creatures []models.Bosses, err error) {
	var allCreatures []models.Bosses
	database.DB.Find(&allCreatures)

	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	return
}

func (ec *EntitysCreatures) GetBosses() (bosses []models.Bosses, err error) {
	if err = database.DB.Find(&bosses).Error; err != nil {
		return
	}
	return
}
func (ec *EntitysCreatures) GetIdBosses(idCreature uint) (creatures models.Bosses, err error) {
	if err = database.DB.Where("id = ?", idCreature).First(&creatures).Error; err != nil {
		return
	}
	return
}

func (ec *EntitysCreatures) CreateMonster(monster models.Creatures) (models.Creatures, error) {
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
func (ec *EntitysCreatures) GetThreeBossesRand() (creatures []models.Bosses, err error) {
	var allCreatures []models.Bosses
	database.DB.Find(&allCreatures)

	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	return
}

func (ec *EntitysCreatures) SaveCreature(monster models.Creatures) (models.Creatures, error) {
	if err := database.DB.Save(&monster).Error; err != nil {
		return monster, err
	}
	return monster, nil
}

func (ec *EntitysCreatures) SaveBoss(monster models.Bosses) (models.Bosses, error) {
	if err := database.DB.Save(&monster).Error; err != nil {
		return monster, err
	}
	return monster, nil
}
func (ec *EntitysCreatures) MaxDamageChanges(creature models.Creatures) (name string, cant uint) {

	if creature.Ice >= cant {
		name = "Ice"
		cant = creature.Ice
	}
	if creature.Physical >= cant {
		name = "Physical"
		cant = creature.Physical
	}
	if creature.Earth >= cant {
		name = "Earth"
		cant = creature.Earth
	}
	if creature.Fire >= cant {
		name = "Fire"
		cant = creature.Fire
	}
	if creature.Death >= cant {
		name = "Death"
		cant = creature.Death
	}
	if creature.Energy >= cant {
		name = "Energy"
		cant = creature.Energy
	}
	if creature.Holy >= cant {
		name = "Holy"
		cant = creature.Holy
	}

	return
}

func (ec *EntitysCreatures) MaxDamageChangesBoss(creature models.Bosses) (name string, cant uint) {

	if creature.Ice >= cant {
		name = "Ice"
		cant = creature.Ice
	}
	if creature.Physical >= cant {
		name = "Physical"
		cant = creature.Physical
	}
	if creature.Earth >= cant {
		name = "Earth"
		cant = creature.Earth
	}
	if creature.Fire >= cant {
		name = "Fire"
		cant = creature.Fire
	}
	if creature.Death >= cant {
		name = "Death"
		cant = creature.Death
	}
	if creature.Energy >= cant {
		name = "Energy"
		cant = creature.Energy
	}
	if creature.Holy >= cant {
		name = "Holy"
		cant = creature.Holy
	}

	return
}
