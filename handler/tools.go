package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
)

type ToolsHandler struct{}

func (th *ToolsHandler) SharedExpHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	template, err := template.ParseFiles(configuration.PATH_WEB_SHARED_LOOT)
	if err != err {
		log.Println(err)
		return
	}
	template.Execute(w, sm)

}

func (th *ToolsHandler) SharedExpProcess(w http.ResponseWriter, r *http.Request) {

	data := r.FormValue("message")
	if data == "" {

		return
	}
	var toolsManager controller.ToolsController

	err := toolsManager.SharedLoot(data, nil)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(map[string]string{"test": "test"})

}

/*

 */
