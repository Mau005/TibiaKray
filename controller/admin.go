package controller

import (
	"fmt"

	"github.com/Mau005/KraynoSerer/components"
	"github.com/Mau005/KraynoSerer/configuration"
)

type AdminController struct{}

func (ac *AdminController) TodaysAproved() string {
	var api ApiController
	var comp components.Components

	var todaysController TodaysController

	todays, _ := todaysController.GetAllTodaysStatus(0)

	title := comp.CreateColsTable("Usuario", "Titulo", "Contenido", "Ver", "Aprobar")
	contentRows := ""
	for _, object := range todays {
		contentRows += comp.CreateRowsTable(
			object.Account.Name,
			object.Title,
			api.NormalizeString(30, object.Description),
			comp.CreateButtonForm("get", fmt.Sprintf(configuration.ROUTER_TODAYS_POST, object.ID), "Ver"),
			comp.CreateButtonForm("post", fmt.Sprintf(configuration.ROUTER_TODAYS_POST_APROVED, object.ID), "Aprobar"))
	}

	return comp.CreateTable(comp.CreateRowsTableFinally(title + contentRows))

}
