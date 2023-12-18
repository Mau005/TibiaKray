package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/Mau005/KraynoSerer/components"
	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/models"
)

type AdminController struct{}

func (ac *AdminController) LobbyAdmin() (result string) {
	var accControl AccountController
	var components components.Components

	countAcc, err := accControl.GetCountAccount()
	if err != nil {
		log.Println(err)
	}
	result += components.CreateLabelADiv(fmt.Sprintf("Cantidad de Usuarios Registrados: %d", countAcc))

	var todaysControll TodaysController

	countTodays, err := todaysControll.GetCountTodays(0)
	if err != nil {
		log.Println(err)
	}
	result += components.CreateLabelADiv(fmt.Sprintf("Cantidad de Todays No Aprovados: %d", countTodays))

	countTodays, err = todaysControll.GetCountTodays(1)
	if err != nil {
		log.Println(err)
	}
	result += components.CreateLabelADiv(fmt.Sprintf("Cantidad de Todays Aprovados: %d", countTodays))

	result += components.CreateLabelADiv(fmt.Sprintf("Cantidad de visitas el dia de hoy: %d", CountVisit))

	var playerCTL PlayerController
	playerCount, err := playerCTL.GetCountPlayer()
	if err != nil {
		log.Println(err)
	}
	result += components.CreateLabelADiv(fmt.Sprintf("Cantidad de Player Registrados: %d", playerCount))
	result += components.CreateLabelADiv(fmt.Sprintf("Hora del Server Save: %d", configuration.Config.Server.ServerSave))
	now := time.Now().Format("2006-01-02 15:04:05")
	result += components.CreateLabelADiv("Hora Actual del servidor: " + now)

	return
}

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

func (ac *AdminController) Streamers() string {
	var api ApiController
	var comp components.Components

	var streamersController StreamerController

	streamer, _ := streamersController.GetStreamers()

	title := comp.CreateColsTable("Streamer", "Titulo", "Tipo de URL", "URL", "Ver")
	contentRows := ""
	for _, object := range streamer {
		contentRows += comp.CreateRowsTable(
			object.Name,
			api.NormalizeString(30, object.Title),
			object.TypeUrl,
			object.URL,
			comp.CreateButtonForm("get", fmt.Sprintf(configuration.ROUTER_STREAMER_ID, object.ID), "ver"))
	}

	formCreateStream := `
		<div>
			<form action="/auth/streamer" method="POST">
				<label for="nombre">Nombre:</label>
				<input type="text" id="nombre" name="nombre" required>

				<label for="titulo">Título:</label>
				<textarea id="titulo" name="titulo" required></textarea>

				<label for="url">Nombre del canal:</label>
				<input type="text" id="url" name="url" required>

				<button type="submit">Enviar</button>
			</form>
		<div>
		`

	return comp.CreateTable(comp.CreateRowsTableFinally(title+contentRows)) + formCreateStream
}

func (ac *AdminController) UserRegister() string {
	var comp components.Components
	title := comp.CreateColsTable("Usuario", "Correo", "Acceso", "Fecha Creación", "Acción")
	var accControl AccountController
	accounts, err := accControl.GetAccountAll()
	if err != nil {
		log.Println(err)
		return ""
	}
	contentRows := ""
	for _, value := range accounts {
		contentRows += comp.CreateRowsTable(
			value.Name,
			value.Email,
			fmt.Sprintf("%d", value.Access),
			value.CreatedAt.Format("2006-01-02 15:04:05"),
			comp.CreateButtonForm("POST",
				fmt.Sprintf("/auth/delete_account/%d", value.ID), "Eliminar"))
	}

	return comp.CreateTable(comp.CreateRowsTableFinally(title + contentRows))

}

func (ac *AdminController) StreamerViews(stream models.Streamers) (content string) {
	content = fmt.Sprintf(`
	<div>
		<form action="/auth/streamer_update" method="POST">
			<input type="hidden" name="id" value="%d" required>
			<label for="nombre">Nombre:</label>
			<input type="text" id="nombre" name="nombre" value="%s" required>

			<label for="titulo">Título:</label>
			<textarea id="titulo" name="titulo" required>%s</textarea>

			<label for="url">Nombre del canal:</label>
			<input type="text" id="url" name="url" value="%s" required>

			<button type="submit">Enviar</button>
		</form>
	<div>
	`, stream.ID, stream.Name, stream.Title, stream.URL)

	return
}
