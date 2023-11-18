package handler

import (
	"net/http"
	"text/template"

	"github.com/Mau005/KraynoSerer/controller"
)

func Page404(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	cs := api.GetBaseWeb(r)

	template, err := template.ParseFiles("static/error404.html")
	if err != nil {
		return
	}

	cs.TitleError = "Un Mensaje de Prueba"
	cs.MSGError = "Esto simula ser un mensaje \n \n tteashkjsdahkjsdhjksd"
	cs.NameButtonError = "Volver a tu hermana"
	cs.RouterError = "/"
	template.Execute(w, cs)

}
