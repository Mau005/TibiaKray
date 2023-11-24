package handler

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
)

type ErrorHandler struct{}

// Indirect page

func (eh *ErrorHandler) DefaultError(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sc := api.GetBaseWeb(r)

	template, err := template.ParseFiles(configuration.PATH_WEB_ERROR)
	if err != nil {
		return
	}
	sc.NameButtonError = "Home"
	sc.MSGError = "Codigo de error 404, esta web esta en mantenimiento o se encuentra en contrucción"
	sc.TitleError = fmt.Sprintf("Opps! Error %d: %s", http.StatusUnauthorized, "Contraseñas incorrectas")
	sc.RouterError = configuration.ROUTER_INDEX
	template.Execute(w, sc)
}

func (eh *ErrorHandler) PageErrorMSG(codeerror, codeMSG int, router string, w http.ResponseWriter, r *http.Request, sc models.StructModel) {

	template, err := template.ParseFiles(configuration.PATH_WEB_ERROR)
	if err != nil {
		return
	}
	sc.TitleError = fmt.Sprintf("Opps! %s: %d!", controller.LenguajeInternal[sc.LenguajeDefault][configuration.ErrorCode], codeerror)
	sc.MSGError = controller.LenguajeInternal[sc.LenguajeDefault][codeMSG]
	if router == "" {
		sc.RouterError = configuration.ROUTER_INDEX
	} else {
		sc.RouterError = router
	}
	sc.NameButtonError = "Return"

	template.Execute(w, sc)

}

func (eh *ErrorHandler) PageErrorStructModel(w http.ResponseWriter, r *http.Request, sc models.StructModel) {

	template, err := template.ParseFiles(configuration.PATH_WEB_ERROR)
	if err != nil {
		return
	}

	template.Execute(w, sc)

}
