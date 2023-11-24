package controller

import (
	"errors"
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

func (cc *CollectorController) GetPlayer() (pl models.Player, err error) {
	words := cc.getColletor("table", configuration.TIBIA_CHARS)
	pl, err = cc.procesingPlayers(words[0])
	if err != nil {
		return pl, err
	}
	return pl, nil
}

func (cc *CollectorController) procesingPlayers(content string) (pl models.Player, err error) {
	keys := []string{"Name:", "Title:", "Level:", "Achievement", "Sex:", "Vocation:", "Points:", "World:", "Residence:", "Guild\u00a0Membership:", "Guild Membership:", "Login:", "CETComment:", "Account Status:", "Account\u00a0Status:"}
	if content == "" {
		return pl, errors.New(configuration.COLLECTOR_EMPTY)
	}
	for _, value := range keys {
		content = strings.ReplaceAll(content, value, "&&")
	}
	formated := strings.Split(content, "&&")
	normalize := []string{}

	if !(len(formated) > 10) {
		return pl, errors.New(configuration.COLLECTOR_NOT_COMPLETED)
	}

	for _, value := range formated {
		if value[0] == 32 || string(value[0]) == " " {
			continue
		}
		normalize = append(normalize, value)
	}

	pl.Name = normalize[0]
	pl.Title = normalize[1]
	pl.Sex = normalize[2]
	pl.Vocation = normalize[3]
	pl.Level = normalize[4]
	pl.Achievement = normalize[5]
	pl.World = normalize[6]
	pl.Residence = normalize[7]
	pl.Guild = normalize[8]
	pl.Login = normalize[9]
	pl.CETComment = normalize[10]
	pl.AccountStatus = normalize[11]
	return pl, nil
}
