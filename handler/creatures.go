package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
)

type CreaturesHandler struct{}

func (ch *CreaturesHandler) CreaturesHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var creaturesController controller.EntitysCreatures
	creatures, err := creaturesController.GetCreatures()

	structNew := struct {
		models.StructModel
		Creatures []models.Creatures
	}{
		StructModel: sm,
		Creatures:   creatures,
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_CREATURES)
	if err != nil {
		log.Println(err)
		return
	}

	templ.Execute(w, structNew)
}
func (ch *CreaturesHandler) BossesHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var creaturesController controller.EntitysCreatures
	Bosses, err := creaturesController.GetBosses()
	structNew := struct {
		models.StructModel
		Bosses []models.Bosses
	}{
		StructModel: sm,
		Bosses:      Bosses,
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_BOSSES)
	if err != nil {
		log.Println(err)
		return
	}

	templ.Execute(w, structNew)
}
