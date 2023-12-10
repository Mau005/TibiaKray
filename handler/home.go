package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gorilla/mux"
)

type HomeHandler struct{}

func (hh *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	controller.CountVisit += 1

	var todaysManager controller.TodaysController

	data, err := todaysManager.GetTodaysLobby()
	if err != nil {
		log.Println(err)
		return
	}

	templ, err := template.ParseFiles("static/index.html")
	if err != nil {
		return
	}
	var api controller.ApiController
	sc := api.GetBaseWeb(r)

	modelNew := struct {
		models.StructModel
		Todays          []models.Todays
		SharedLootHight models.SharedLoot
		Rashid          string
	}{
		StructModel:     sc,
		Todays:          data,
		SharedLootHight: configuration.SharedLootHightNow,
		Rashid:          configuration.Rashid,
	}

	templ.Execute(w, modelNew)

}

func (hh *HomeHandler) Todays(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("static/todays.html")
	if err != nil {
		return
	}
	var api controller.ApiController
	sc := api.GetBaseWeb(r)
	var todayCont controller.TodaysController

	todays, _ := todayCont.GetTodayView()

	todaysWeb := struct {
		models.StructModel
		Todays []models.Todays
	}{
		StructModel: sc,
		Todays:      todays,
	}
	templ.Execute(w, todaysWeb)

}

func (hh *HomeHandler) TodaysPost(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sc := api.GetBaseWeb(r)
	arg := mux.Vars(r)

	var ErrorHandler ErrorHandler

	idTodays, err := strconv.ParseUint(arg["id"], 10, 64)
	if err != nil {
		ErrorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorInternal, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, idTodays), w, r, sc)
		return
	}

	var todaysController controller.TodaysController

	todays, err := todaysController.GetToday(uint(idTodays))
	if err != nil {
		ErrorHandler.PageErrorMSG(http.StatusNotFound, configuration.ErrorDefault, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, idTodays), w, r, sc)
		return
	}

	templ, err := template.ParseFiles("static/todays_post.html")
	if err != nil {
		ErrorHandler.PageErrorMSG(http.StatusNotFound, configuration.ErrorDefault, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, idTodays), w, r, sc)
		return
	}

	var accManager controller.AccountController
	var votedManager controller.VotedController

	var voted models.Voted
	account, err := accManager.GetAccount(sc.Email)
	//Iniciamos el proceso de buscar si el usuario tiene un voto en la web
	if err == nil {
		voted, _ = votedManager.GetVotedAccountTodays(todays.ID, account.ID) //inicio el voto si existe en la cuenta
	}

	todaysWeb := struct {
		models.StructModel
		Todays models.Todays
		Voted  models.Voted
	}{
		StructModel: sc,
		Todays:      todays,
		Voted:       voted,
	}
	templ.Execute(w, todaysWeb)
}
