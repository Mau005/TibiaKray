package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
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
