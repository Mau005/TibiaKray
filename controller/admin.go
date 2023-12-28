package controller

import (
	"fmt"
	"log"
	"strings"
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

func (ac *AdminController) NewsTicket() string {
	var api ApiController
	var comp components.Components

	var newsTicketController NewsTicketController

	newsTicket, _ := newsTicketController.GetAllNewsTicket()

	title := comp.CreateColsTable("Categoria", "Esp", "Fecha Creacion", "URL")
	contentRows := ""
	for _, object := range newsTicket {
		contentRows += comp.CreateRowsTable(
			object.Title,
			api.NormalizeString(30, object.Content),
			object.CreatedAt.Format("2006-01-02 15:04:05"),
			comp.CreateButtonForm("get", fmt.Sprintf(configuration.ROUTER_NEWS_TICKET_ID, object.ID), "ver"))
	}
	return comp.CreateTable(comp.CreateRowsTableFinally(title + contentRows))
}

func (ac *AdminController) ViewNewsTicket(ticket models.NewsTicket, type_form string) string {
	checket := `<input type="checkbox" name="StatusNews"><br>`
	if ticket.StatusNews {
		checket = `<input type="checkbox" name="StatusNews" checked><br>`
	}
	select_value := ""
	if !(ticket.Title == "") {
		select_value = fmt.Sprintf(`<option value="%s">%s</option>`, ticket.Title, strings.ToUpper(ticket.Title[:1])+ticket.Title[1:])
	}

	return fmt.Sprintf(`
	<div>
		<form action="/auth/newsticket" method="post">
			<input type="hidden" name="id_ticket" value="%d"><br>
			<input type="hidden" name="typeForm" value="%s"><br>

			<label for="title">Title:</label>
			<select class="form-control" name="Title">
					%s
					<option value="update">Update</option>
					<option value="internal">Internal</option>
					<option value="tibia">Tibia</option>
				</select>

			<label for="content">Content:</label>
			<textarea id="content" name="Content" required>%s</textarea><br>

			<label for="contentEn">Content (English):</label>
			<textarea id="contentEn" name="ContentEn">%s</textarea><br>

			<label for="contentPl">Content (Polish):</label>
			<textarea id="contentPl" name="ContentPl">%s</textarea><br>

			<label for="contentBr">Content (Portuguese):</label>
			<textarea id="contentBr" name="ContentBr">%s</textarea><br>

			<label for="statusNews">Status News:</label>
			%s

			<button type="submit">Submit</button>
		</form>
	<div>	
	`, ticket.ID, type_form, select_value, ticket.Content, ticket.ContentEn, ticket.ContentPl, ticket.ContentBr, checket)
}
