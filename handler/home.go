package handler

import (
	"net/http"
	"text/template"

	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
)

type HomeHandler struct{}

func (hh *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("static/index.html")
	if err != nil {
		return
	}
	var api controller.ApiController
	sc := api.GetBaseWeb(r)
	templ.Execute(w, sc)

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
