package handler

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gorilla/mux"
)

type AdminHandler struct{}

func (a *AdminHandler) Lobby(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	if !(sm.Access > 2) {
		return
	}

	var adminCTL controller.AdminController
	contentAdmin := adminCTL.LobbyAdmin()

	template, err := template.ParseFiles(configuration.PATH_WEB_ADMIN)
	if err != nil {
		return
	}

	content := struct {
		models.StructModel
		Title   string
		Content string
	}{
		StructModel: sm,
		Title:       "Información de la Aplicación",
		Content:     contentAdmin,
	}

	template.Execute(w, content)

}

func (a *AdminHandler) TodaysAproved(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var admin controller.AdminController

	if !(sm.Access > 2) {
		return
	}

	template, err := template.ParseFiles("static/admin.html")
	if err != nil {
		return
	}

	content := struct {
		models.StructModel
		Title   string
		Content string
	}{
		StructModel: sm,
		Title:       "Panel de Adminsitracion de Aprovaciones de HOY",
		Content:     admin.TodaysAproved(),
	}

	template.Execute(w, content)

}

func (a *AdminHandler) TodaysAprovedPOST(w http.ResponseWriter, r *http.Request) {
	arg := mux.Vars(r)

	idTodays, err := strconv.ParseUint(arg["id"], 10, 64)
	if err != nil {
		return
	}

	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var accManager controller.AccountController
	acc, err := accManager.GetAccount(sm.Email)
	if err != nil {
		return
	}

	if !(acc.Access >= configuration.Config.Access.AprovedTodays) {
		return
	}

	var todaysManager controller.TodaysController

	todays, err := todaysManager.GetToday(uint(idTodays))
	if err != nil {
		return
	}
	todays.Status = true
	_, err = todaysManager.SaveTodays(todays)
	if err != nil {
		return
	}

	var admin controller.AdminController

	if !(sm.Access > 2) {
		return
	}

	template, err := template.ParseFiles("static/admin.html")
	if err != nil {
		return
	}

	content := struct {
		models.StructModel
		Title   string
		Content string
	}{
		StructModel: sm,
		Title:       "Panel de Adminsitracion de Aprovaciones de HOY",
		Content:     admin.TodaysAproved(),
	}

	template.Execute(w, content)

}

func (a *AdminHandler) UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	if !(sm.Access >= configuration.ACCES_ADMIN) {
		var errorHandler ErrorHandler
		errorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, configuration.ROUTER_INDEX, w, r, sm)
		return
	}

	var adminCtl controller.AdminController

	strucNew := struct {
		models.StructModel
		Title   string
		Content string
	}{
		StructModel: sm,
		Title:       "Usuarios Registrados",
		Content:     adminCtl.UserRegister()}

	template, err := template.ParseFiles(configuration.PATH_WEB_ADMIN)
	if err != nil {
		log.Println(err)
	}
	template.Execute(w, strucNew)
}
