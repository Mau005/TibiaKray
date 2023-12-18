package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gorilla/mux"
)

type RecoveryHandler struct{}

func (rh *RecoveryHandler) RecoveryHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)
	var errorHandler ErrorHandler

	arg := mux.Vars(r)

	if arg["code"] == "" {
		var errorHandler ErrorHandler
		errorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorEmptyField, configuration.ROUTER_INDEX, w, r, sm)
		return
	}

	ra, ok := controller.RecoveryAccount[arg["code"]]
	if !ok {
		errorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, "", w, r, sm)
		return
	}

	if time.Now().After(ra.ExpireAt) {
		errorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorExpiredProcess, "", w, r, sm)
		return
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_RECOVERY)
	if err != nil {
		log.Println(err)
		return
	}

	structNew := struct {
		models.StructModel
		Code string
	}{
		StructModel: sm,
		Code:        arg["code"],
	}
	templ.Execute(w, structNew)
}

func (rh *RecoveryHandler) RecoveryAccount(w http.ResponseWriter, r *http.Request) {
	var errorHandler ErrorHandler
	var api controller.ApiController
	sm := api.GetBaseWeb(r)
	nameOREmail := r.FormValue("recovery")
	if nameOREmail == "" {
		errorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorEmptyField, "", w, r, sm)
		return

	}

	var accController controller.AccountController

	acc, err := accController.GetAccount(nameOREmail)
	if err != nil {
		acc, err = accController.GetAccountUser(nameOREmail)
		if err != nil {
			errorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, "", w, r, sm)
			return
		}
	}

	var recoveryController controller.RecoveryController
	rc, err := recoveryController.NewRecoveryAccount(acc.Name, acc.Email, acc.Languaje)
	if err != nil {
		log.Println(err)
	}
	controller.RecoveryAccount[rc.Hash] = rc

	var emailController controller.EmailController
	url := fmt.Sprintf(configuration.ROUTER_RECOVERY_CODE, configuration.Config.Server.Ip, rc.Hash)

	file, err := os.ReadFile(configuration.PATH_RECOVERY_EMAIL)
	if err != nil {
		log.Println(err)
		return
	}

	container := emailController.GenerateEmailSend(
		controller.LenguajeInternal[sm.LenguajeDefault][configuration.AccountRecovery],
		fmt.Sprintf(string(file), acc.Name, url, rc.Code))
	go emailController.SendEmail(acc.Email, container)
	errorHandler.PageMSG(configuration.ProcessedSuccess, configuration.ChangePassword, "", w, r, sm)
}

func (rh *RecoveryHandler) RecoveryChangePassword(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("hash")
	code := r.FormValue("code")
	password := r.FormValue("recoverypassword")
	passwordTwo := r.FormValue("recoverypasswordtwo")

	var api controller.ApiController
	var errorHandler ErrorHandler

	sm := api.GetBaseWeb(r)

	//controlo si los password son necesarios
	if !(len(password) >= configuration.MAX_LEN_PASSWORD && password == passwordTwo) {
		errorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorEmptyField, "", w, r, sm)
		return
	}

	rc, ok := controller.RecoveryAccount[hash]
	if !ok {
		errorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorCode, "", w, r, sm)
		return
	}

	if !(rc.Code == code) {
		errorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorCode, "", w, r, sm)
		return
	}

	var accController controller.AccountController
	acc, err := accController.GetAccount(rc.Email)
	if err != nil {
		errorHandler.PageErrorMSG(http.StatusInternalServerError, configuration.ErrorInternal, "", w, r, sm)
		return
	}

	acc.Password = api.GenerateCryptPassword(password)
	accController.SaveAccount(acc)

	errorHandler.PageMSG(configuration.ProcessedSuccess, configuration.SavedSuccess, "", w, r, sm)

}
