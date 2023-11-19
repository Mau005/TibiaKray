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

	if email == "" || password == "" {
		return
	}

	var accController controller.AccountController

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
	fmt.Println("Inicio Session")

}

func (ac *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	policy := r.FormValue("policy")
	username := r.FormValue("username")
	passwordTwo := r.FormValue("passwordTwo")
	password := r.FormValue("password")

	if policy != "on" {
		return
	}

	if password == "" {
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
	acc, err := accController.CreateAccount(email, username, password)
	if err != nil {
		json.NewEncoder(w).Encode(models.Exception{
			Error:     configuration.ERROR_SERVICE_ACCOUNT,
			Status:    http.StatusNetworkAuthenticationRequired,
			Message:   err.Error(),
			TimeStamp: time.Now(),
		})
		return
	}
	json.NewEncoder(w).Encode(acc)

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
