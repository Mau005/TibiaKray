package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
)

type ToolsHandler struct{}

type SharedLootController struct {
	models.StructModel
	SharedLootHight models.SharedLoot
	ProcesingData   map[string][]string
}

func (th *ToolsHandler) SharedLootHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	template, err := template.ParseFiles(configuration.PATH_WEB_SHARED_LOOT)
	if err != err {
		log.Println(err)
		return
	}
	template.Execute(w, SharedLootController{StructModel: sm, SharedLootHight: configuration.SharedLootHightNow})
}

func (th *ToolsHandler) SharedLootProcess(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("message")
	var api controller.ApiController

	sm := api.GetBaseWeb(r)
	template, err := template.ParseFiles(configuration.PATH_WEB_SHARED_LOOT)
	if err != err {
		log.Println(err)
		return
	}

	if data == "" {
		sm.MSGError = "Protocolo No Sportado"
		template.Execute(w, SharedLootController{StructModel: sm, SharedLootHight: configuration.SharedLootHightNow})
		return
	}
	var toolsManager controller.ToolsController

	_, sharedMap, err := toolsManager.SharedLoot(data)
	if err != nil {
		sm.MSGError = err.Error()
		template.Execute(w, SharedLootController{StructModel: sm, SharedLootHight: configuration.SharedLootHightNow})
		return
	}

	template.Execute(w, SharedLootController{StructModel: sm, SharedLootHight: configuration.SharedLootHightNow, ProcesingData: sharedMap})

}

/*

 */

func (th *ToolsHandler) ToolsHandlerItems(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	template, err := template.ParseFiles(configuration.PATH_WEB_TOOLS)
	if err != nil {
		log.Println(err)
	}
	template.Execute(w, sm)
}
