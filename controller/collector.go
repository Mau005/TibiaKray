package controller

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gocolly/colly/v2"
)

var News *[]models.News
var NewsTicket *[]models.NewsTicket

type CollectorController struct{}

func (cc *CollectorController) GenerateNewsTibia() (err error) {
	news, err := cc.GetNews()
	if err != nil {
		log.Println(err)
	}

	News = &news

	newsT, err := cc.GetNewsTickets()
	if err != nil {
		return err
	}
	NewsTicket = &newsT

	return nil
}

func (cc *CollectorController) GetNews() (t []models.News, err error) {
	key := []string{"\u00a0", "\u00a0"}
	words := cc.getColletor("div.Border_2", configuration.TIBIA_NEWS)
	content := words[1]
	for _, value := range key {
		content = strings.ReplaceAll(content, value, " ")
	}

	listContent := strings.Split(content, "\n")
	normalize := []string{}
	for _, value := range listContent {
		if value == "" {
			continue
		}
		normalize = append(normalize, value)
	}

	t, err = cc.procesingNews(normalize)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (cc *CollectorController) procesingNews(content []string) (newsT []models.News, err error) {

	keys := map[string]bool{
		"Jan": true,
		"Feb": true,
		"Mar": true,
		"Apr": true,
		"May": true,
		"Jun": true,
		"Jul": true,
		"Aug": true,
		"Sep": true,
		"Oct": true,
		"Nov": true,
		"Dec": true,
	}

	groupIndex := []int{}
	for index, value := range content {
		if len(value) > 2 {
			if keys[value[:3]] {
				groupIndex = append(groupIndex, index)
			}
		}
	}
	groupIndex = append(groupIndex, -1)
	if !(len(groupIndex)%2 == 0) {
		return newsT, errors.New("No Soporta el Protocolo de converison")
	}

	group := make(map[int][]string)
	aux := 0
	for i := 0; i < len(groupIndex); i += 2 {
		aux++
		if i == 0 {
			group[aux] = content[i:groupIndex[i+2]]
		} else {
			group[aux] = content[groupIndex[i-2]:groupIndex[i]]
		}
	}

	newsT, err = cc.subProcesingNews(group)
	if err != nil {
		return newsT, err
	}

	return newsT, err
}

func (cc *CollectorController) subProcesingNews(group map[int][]string) (newsT []models.News, err error) {

	count := len(group)
	for i := 1; i <= count; i++ {
		prosc := group[i]
		title := cc.calculatorString(false, prosc[:2])
		content := cc.calculatorString(true, prosc[2:])
		newsT = append(newsT, models.News{
			Title:     title,
			Content:   content,
			AccountID: nil,
		})
	}
	return newsT, err
}

func (cc *CollectorController) calculatorString(lineJumper bool, content []string) (result string) {
	if !(lineJumper) {
		for _, value := range content {
			result += value
		}
	} else {
		for _, value := range content {
			result += value + "\n"
		}
	}
	return result
}
func (cc *CollectorController) GetNewsTickets() (news []models.NewsTicket, err error) {
	key := []string{"\u00a0", "\u00a0"}
	words := cc.getColletor("div.BoxContent", configuration.TIBIA_NEWS)

	content := words[0]
	for _, value := range key {
		content = strings.ReplaceAll(content, value, " ")
	}
	listContent := strings.Split(content, "\n")
	normalize := []string{}
	for _, value := range listContent {
		if value == "" {
			continue
		}
		normalize = append(normalize, value)
	}
	t, err := cc.procesinNewsTicket(normalize)

	return t, nil
}

func (cc *CollectorController) procesinNewsTicket(content []string) (newsT []models.NewsTicket, err error) {
	if len(content) == 0 {
		return newsT, errors.New("Not DATA")
	}
	atrContent := len(content) % 4

	if !(atrContent == 0) {
		return newsT, errors.New("Protocol Incompatible NewsTickets")
	}

	for i := 0; i < len(content); i += 4 {
		newsT = append(newsT, models.NewsTicket{
			Title:     content[i] + " " + content[i+1],
			Content:   content[i+3],
			AccountID: nil,
		})

	}

	return newsT, err
}

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
