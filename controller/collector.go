package controller

import (
	"fmt"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gocolly/colly/v2"
)

var News *[]models.News
var NewsTicket *[]models.NewsTicket

type CollectorController struct{}

func (cc *CollectorController) getColletor(keyXML, site string) (dat []string) {
	c := colly.NewCollector()
	words := []string{}

	c.OnHTML(keyXML, func(e *colly.HTMLElement) {
		words = append(words, e.Text)
	})
	c.Visit(site)

	return words
}

func (cc *CollectorController) GetPlayer(name_character string) (pl models.Player, err error) {
	name_character = strings.Trim(name_character, " ")
	name_character = strings.ReplaceAll(name_character, " ", "+")
	words := cc.getColletor("table.TableContent", fmt.Sprintf(configuration.TIBIA_CHARS, name_character))

	pl, err = cc.procesingPlayers(words[0])

	return pl, nil
}

func (cc *CollectorController) procesingPlayers(response string) (pl models.Player, err error) {

	// Lista de claves
	//keys := []string{"Name:", "Title:", "Level:", "Achievement", "Sex:", "Vocation:", "Points:", "World:", "Residence:", "Guild\u00a0Membership:", "Guild Membership:", "Login:", "CETComment:", "Account Status:", "Account\u00a0Status:"}

	// Crear la expresión regular
	/*
		re := regexp.MustCompile(`(?i)` + // Ignorar mayúsculas y minúsculas
			`(` +                       // Inicio de la captura
			`Name:|Title:|Level:|Achievement|Sex:|Vocation:|Points:|World:|Residence:|Guild\s?Membership:|Login:|CETComment:|Account\s?Status:|Account\u00a0Status:` + // Alternancia de claves
			`)` + // Fin de la captura
			`([^A-Z]+)` // Captura de cualquier cosa que no sea letras mayúsculas
	*/
	//fmt.Printf("%q", response)
	response = strings.ReplaceAll(response, "\u00a0", " ")
	//fmt.Printf("%q", response)

	keys := []string{"", //name por referencia mas abajo
		"Former Names",
		"Title",
		"Sex",
		"Vocation",
		"Level",
		"Achievement Points",
		"World",
		"Former World",
		"Residence",
		"Comment",
		"House",
		"Guild Membership",
		"Last Login",
		"Account Status",
	}
	playerData := make(map[string]string, len(keys))
	for _, keysData := range keys {
		playerData[keysData] = ""
	}

	subProcesing := strings.Split(response, ":")
	KeyOld := "Name"
	for i, value := range subProcesing {
		if i == 0 {
			continue
		}

		for iKeys, keysValues := range keys {
			if len(value) >= len(keysValues) {
				if keysValues == "" {
					continue
				}
				contentAttrResult := value[:len(value)-len(keysValues)]
				keyAttrResult := value[len(value)-len(keysValues):]
				//fmt.Println(value[:len(value)-len(keysValues)]) asi saco el valor
				//value[len(value)-len(keysValues):] saco el contenido

				if keyAttrResult == keysValues {
					//fmt.Printf("Check: K:%s V:%s\n", keyAttrResult, keysValues)
					//fmt.Printf("Asign: K:%s V:%s\n", KeyOld, contentAttrResult)
					playerData[KeyOld] = contentAttrResult
					KeyOld = keyAttrResult
					keys[iKeys] = ""
					break
				}
			}
		}
		if len(subProcesing)-1 == i {
			playerData[KeyOld] = value
		}

	}

	return models.Player{
		Name:          playerData["Name"],
		FormerNames:   playerData["Former Names"],
		Title:         playerData["Title"],
		Sex:           playerData["Sex"],
		Vocation:      playerData["Vocation"],
		Level:         playerData["Level"],
		Achievement:   playerData["Achievement Points"],
		World:         playerData["World"],
		FormerWorld:   playerData["Former World"],
		Residence:     playerData["Residence"],
		Comment:       playerData["Comment"],
		House:         playerData["House"],
		Guild:         playerData["Guild Membership"],
		LastLogin:     playerData["Last Login"],
		AccountStatus: playerData["Account Status"],
	}, nil
}
