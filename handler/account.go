package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
)

type AccountHandler struct{}

func (ac *AccountHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("user")
	password := r.FormValue("passworduser")
	var ApiController controller.ApiController
	sc := ApiController.GetBaseWeb(r)
	var errHandler ErrorHandler

	if email == "" || password == "" {
		errHandler.PageErrorMSG(http.StatusNoContent, configuration.ErrorEmptyField, "", w, r, sc)
		return
	}

	var accController controller.AccountController

	token, err := accController.Login(email, password)
	if err != nil {
		errHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPassword, "", w, r, sc)
		return
	}
	var api controller.ApiController
	api.SaveSession(&token, w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (ac *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	policy := r.FormValue("policy")
	username := r.FormValue("username")
	passwordTwo := r.FormValue("passwordTwo")
	password := r.FormValue("password")

	if policy != "on" {
		var ApiController controller.ApiController
		sc := ApiController.GetBaseWeb(r)
		var errHandler ErrorHandler
		errHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorPolicies, "", w, r, sc)
		return
	}

	if password == "" {
		var ApiController controller.ApiController
		sc := ApiController.GetBaseWeb(r)
		var errHandler ErrorHandler
		errHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorEmptyField, "", w, r, sc)
		return
	}

	if passwordTwo != password {
		return
	}

	var accController controller.AccountController

	conditions := accController.IsUserName(username)
	if conditions {
		return
	}
	_, err := accController.CreateAccount(email, username, password)
	if err != nil {
		json.NewEncoder(w).Encode(models.Exception{
			Error:     configuration.ERROR_SERVICE_ACCOUNT,
			Status:    http.StatusNetworkAuthenticationRequired,
			Message:   err.Error(),
			TimeStamp: time.Now(),
		})
		return
	}
	token, err := accController.Login(email, password)
	if err != nil {
		json.NewEncoder(w).Encode(models.Exception{
			Error:     configuration.ERROR_SERVICE_ACCOUNT,
			Status:    http.StatusNetworkAuthenticationRequired,
			Message:   configuration.ERROR_PRIVILEGES_GEN,
			TimeStamp: time.Now(),
		})
		return
	}
	var api controller.ApiController
	api.SaveSession(&token, w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (ac *AccountHandler) MyProfileHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	cs := api.GetBaseWeb(r)

	if cs.Email == "" {
		json.NewEncoder(w).Encode(models.Exception{Message: "No tienes permiso HIDEPUTA"})
		return
	}

	template, err := template.ParseFiles("static/my_profile.html")
	if err != nil {
		return
	}

	prueba := struct {
		models.StructModel
		Test string
	}{
		StructModel: cs,
		Test:        "hola",
	}

	template.Execute(w, prueba)
}

func (ac *AccountHandler) MyProfileSettingPOST(w http.ResponseWriter, r *http.Request) {
	streamMode := r.FormValue("streammode")
	lenguaje := r.FormValue("lenguaje")

	var api controller.ApiController

	sm := api.GetBaseWeb(r)

	var acc controller.AccountController

	account, err := acc.GetAccount(sm.Email)
	if err != nil {
		fmt.Println(err)
		return
	}

	account.StreamMode = streamMode == "on"

	if lenguaje != "" {
		baseLenguaje := make([]string, 0, len(controller.Lenguaje)-1)
		for key := range controller.Lenguaje {
			if key == "base" {
				continue
			}
			baseLenguaje = append(baseLenguaje, key)
		}
		base := controller.Lenguaje[account.Languaje]
		for _, value := range baseLenguaje {
			if base[value] == lenguaje {
				account.Languaje = value
			}
		}
	}
	err = acc.UpdateSession(account, w, r)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/auth/my_profile", http.StatusSeeOther)
}

func (ac *AccountHandler) MyProfileChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	passwordTwo := r.FormValue("password2")
	if password != passwordTwo || password == "" || passwordTwo == "" {
		return
	}

	var api controller.ApiController
	newPassword := api.GenerateCryptPassword(password)

	acc, err := api.GetSessionAccount(r)
	if err != nil {
		return
	}
	acc.Password = newPassword
	var accController controller.AccountController
	accController.SaveAccount(acc)
	ac.Logout(w, r)

}

func (ac *AccountHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	api.SaveSession(nil, w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ac *AccountHandler) AddCommentTodays(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comments")
	id := r.FormValue("id")

	var api controller.ApiController
	seccion, err := api.GetSessionAccount(r)
	if err != nil {
		return
	}
	var accManager controller.AccountController
	account, err := accManager.GetAccount(seccion.Email)
	if err != nil {
		return
	}

	comn, err := accManager.AddCommentTodays(id, comment, account)
	if err != nil {
		return
	}

	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, *comn.TodaysID), http.StatusSeeOther)
}

func (ac *AccountHandler) MyProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	cm := api.GetBaseWeb(r)

	var votedManager controller.VotedController
	voteds, err := votedManager.MyVotedImage(cm.Email)
	if err != nil {
		fmt.Println("entra en error voted?")
		return
	}

	strucNew := struct {
		models.StructModel
		Voteds []models.Voted
	}{
		StructModel: cm,
		Voteds:      voteds,
	}

	template, err := template.ParseFiles(configuration.PATH_WEB_MY_FAVO_PIC)
	if err != nil {
		return
	}
	template.Execute(w, strucNew)
}

func (ac *AccountHandler) MyProfilePLayers(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	cm := api.GetBaseWeb(r)

	strucNew := struct {
		models.StructModel
		Voteds []models.Voted
	}{
		StructModel: cm,
	}

	template, err := template.ParseFiles(configuration.PATH_WEB_MY_PLAYERS)
	if err != nil {
		return
	}
	template.Execute(w, strucNew)
}

func (ac *AccountHandler) SearchMyPlayer(w http.ResponseWriter, r *http.Request) {

	character := r.FormValue("nameplayer")
	if character == "" {
		log.Println("Error len = 0")
		return
	}

	var collectorAPI controller.CollectorController
	pl, err := collectorAPI.GetPlayer(character)
	if err != nil {
		log.Println("Error collector")
		return
	}
	json.NewEncoder(w).Encode(pl)
}
