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

func (a *AdminHandler) StreamerHandlerAdmin(w http.ResponseWriter, r *http.Request) {
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
		Title:       "Streamers",
		Content:     adminCtl.Streamers()}

	template, err := template.ParseFiles(configuration.PATH_WEB_ADMIN)
	if err != nil {
		log.Println(err)
	}
	template.Execute(w, strucNew)
}

func (a *AdminHandler) StreamerPOSTAdmin(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	if !(sm.Access >= configuration.ACCES_ADMIN) {
		var errorHandler ErrorHandler
		errorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, configuration.ROUTER_INDEX, w, r, sm)
		return
	}

	name := r.FormValue("nombre")
	title := r.FormValue("titulo")
	url := r.FormValue("url")

	var streamController controller.StreamerController
	_, err := streamController.CreateStremer(models.Streamers{Name: name, Title: title, URL: fmt.Sprintf(configuration.TWITCH_CANAL, url, configuration.Config.Server.Ip)})
	if err != nil {
		log.Println(err)
	}
	a.StreamerHandlerAdmin(w, r)
}

func (a *AdminHandler) StreamerIDAdmin(w http.ResponseWriter, r *http.Request) {

	args := mux.Vars(r)
	idStream, err := strconv.Atoi(args["id"])
	if err != nil {
		log.Println(err)
		return
	}
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var streamController controller.StreamerController

	stream, err := streamController.GetIdStreamers(uint(idStream))
	if err != nil {
		log.Println(err)
		return
	}

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
		Title:       "Streamers",
		Content:     adminCtl.StreamerViews(stream)}

	template, err := template.ParseFiles(configuration.PATH_WEB_ADMIN)
	if err != nil {
		log.Println(err)
	}
	template.Execute(w, strucNew)
}

func (a *AdminHandler) StreamPatchAdmin(w http.ResponseWriter, r *http.Request) {
	idStreamForm := r.FormValue("id")
	name := r.FormValue("nombre")
	title := r.FormValue("titulo")
	url := r.FormValue("url")

	idStream, err := strconv.Atoi(idStreamForm)
	if err != nil {
		log.Println(err)
		return
	}
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var streamController controller.StreamerController

	stream, err := streamController.GetIdStreamers(uint(idStream))
	if err != nil {
		log.Println(err)
		return
	}

	if !(sm.Access >= configuration.ACCES_ADMIN) {
		var errorHandler ErrorHandler
		errorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, configuration.ROUTER_INDEX, w, r, sm)
		return
	}
	stream.Name = name
	stream.URL = url
	stream.Title = title
	_, err = streamController.SaveStremer(stream)
	if err != nil {
		log.Println(err)
		return
	}
	a.StreamerHandlerAdmin(w, r)
}

func (a *AdminHandler) NewsTicketHandler(w http.ResponseWriter, r *http.Request) {
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
		Title:       "NewsTickets",
		Content:     adminCtl.NewsTicket() + adminCtl.ViewNewsTicket(models.NewsTicket{}, "create")}

	template, err := template.ParseFiles(configuration.PATH_WEB_ADMIN)
	if err != nil {
		log.Println(err)
	}
	template.Execute(w, strucNew)
}

func (a *AdminHandler) NewsTicketProcesing(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	if !(sm.Access >= configuration.ACCES_ADMIN) {
		var errorHandler ErrorHandler
		errorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, configuration.ROUTER_INDEX, w, r, sm)
		return
	}

	var ticketController controller.NewsTicketController

	idStr, err := strconv.ParseUint(r.FormValue("id_ticket"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	newsTicket := models.NewsTicket{}
	typeForm := r.FormValue("typeForm")

	if typeForm == "update" {
		newsTicket, err = ticketController.GetIDNewsTicket(uint(idStr))
		if err != nil {
			log.Println(err)
			return
		}
	}
	newsTicket.Title = r.FormValue("Title")
	newsTicket.Content = r.FormValue("Content")
	newsTicket.ContentBr = r.FormValue("ContentBr")
	newsTicket.ContentEn = r.FormValue("ContentEn")
	newsTicket.ContentPl = r.FormValue("ContentPl")
	newsTicket.StatusNews = r.FormValue("StatusNews") == "on"

	if typeForm == "create" {
		var accountController controller.AccountController
		account, err := accountController.GetAccount(sm.Email)
		if err != nil {
			log.Println(err)
			return
		}
		newsTicket.Account = *account
		newsTicket.AccountID = &account.ID
		newsTicket, err = ticketController.CreateNewsTicket(newsTicket)
		if err != nil {
			log.Println(err)
		}
	} else if typeForm == "update" {
		newsTicket, err = ticketController.SaveNewsTicket(newsTicket)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, configuration.ROUTER_NEWS_TICKET, http.StatusSeeOther)
		return
	} else {
		log.Println("No pudo gestionar si actualizar o crear")
	}
	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_NEWS_TICKET_ID, newsTicket.ID), http.StatusSeeOther)
}

func (a *AdminHandler) NewsTicketIDHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	if !(sm.Access >= configuration.ACCES_ADMIN) {
		var errorHandler ErrorHandler
		errorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, configuration.ROUTER_INDEX, w, r, sm)
		return
	}

	var adminCtl controller.AdminController
	var newsTicketController controller.NewsTicketController

	args := mux.Vars(r)

	idTicket, err := strconv.ParseUint(args["id"], 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	newsTicket, err := newsTicketController.GetIDNewsTicket(uint(idTicket))
	if err != nil {
		log.Println(err)
		return
	}
	strucNew := struct {
		models.StructModel
		Title   string
		Content string
	}{
		StructModel: sm,
		Title:       "NewsTickets",
		Content:     adminCtl.ViewNewsTicket(newsTicket, "update")}

	template, err := template.ParseFiles(configuration.PATH_WEB_ADMIN)
	if err != nil {
		log.Println(err)
	}
	template.Execute(w, strucNew)
}
