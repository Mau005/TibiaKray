package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gorilla/mux"
)

type ImageHandler struct{}

// StatusIMprove:
var UploadSecuirtyPass = map[string]bool{
	"Default": true,
	"Image":   true,
	"Mp4":     true,
	"Twitch":  true,
	"Youtube": true,
}

func (su *ImageHandler) UploadHandlerDefault(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles(configuration.PATH_WEB_UPLOAD_FILES)
	if err != nil {
		log.Println(err)
		return
	}
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	structNew := struct {
		models.StructModel
		StatusWeb string
	}{
		StructModel: sm,
		StatusWeb:   "Default",
	}

	templ.Execute(w, structNew)
}
func (su *ImageHandler) UploadDessioningHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	statusWeb := r.FormValue("StatusWeb")
	var errorHandler ErrorHandler

	if statusWeb == "" {
		errorHandler.PageErrorMSG(http.StatusNotFound, configuration.ErrorDefault, configuration.ROUTER_UPLOAD_IMAGES, w, r, sm)
		return
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_UPLOAD_FILES)
	if err != nil {
		log.Println(err)
		return
	}

	structNew := struct {
		models.StructModel
		StatusWeb string
	}{
		StructModel: sm,
		StatusWeb:   statusWeb,
	}

	templ.Execute(w, structNew)
}

func (su *ImageHandler) LoadImage(w http.ResponseWriter, r *http.Request) {

	var api controller.ApiController
	//var accCont controller.AccountController
	var ErrorHandler ErrorHandler

	acc, err := api.GetSessionAccount(r)
	if err != nil {
		//Problem Security solicitude log in web
		//var errorHandler ErrorHandler
		//errorHandler.DefaultError()
		return
	}

	var extencion = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"bmp":  true,
		"mp4":  true,
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	file, handler, err := r.FormFile("documents")
	if err != nil {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}
	defer file.Close()

	filesize := handler.Size
	fmt.Println(filesize, configuration.MAX_FILE_SIZE)
	if filesize >= configuration.MAX_FILE_SIZE {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorMaxFileSize, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}

	verify := strings.Split(handler.Filename, ".")
	extencionFile := verify[len(verify)-1]

	if !(extencion[extencionFile]) {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorPolicies, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}

	var todaysController controller.TodaysController

	todays, err := todaysController.CreateTodays(models.Todays{
		Title:       title,
		Description: description,
		Status:      false,
		Account:     *acc,
		AccountID:   acc.ID,
	})

	if todays.ID == 0 {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}

	newName := api.GenerateHash(fmt.Sprintf("%d", todays.ID)) //Genero un sha256, proceso el id para ese hash, luego sumo la extencon del file
	fileExtencion := newName + "." + extencionFile
	pathSave := fmt.Sprintf("%s/%s/", configuration.IMAGEN_PATH, api.GenerateHash(acc.Name))
	pathConsume := fmt.Sprintf("%s/%s/%s", "todays", api.GenerateHash(acc.Name), fileExtencion)

	err = os.MkdirAll(pathSave, os.ModePerm)
	if err != nil {
		//En esta parte debo eliminar el todays generado
		return
	}

	fileCopy, err := os.Create(pathSave + fileExtencion)
	if err != nil {
		//En esta parte debo eliminar el todays generado
		return
	}

	defer fileCopy.Close()

	_, err = io.Copy(fileCopy, file)
	if err != nil {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}

	var encr models.Files
	encr.TodaysID = &todays.ID
	encr.PathConsume = pathConsume
	encr.PathEncrypt = fileExtencion
	encr.PathOrigin = pathSave + fileExtencion
	encr.Todays = &todays
	if extencionFile == "mp4" {
		encr.TypeFile = "mp4"
	}

	var encrController controller.FileController

	enc, err := encrController.CreateFile(encr)
	if err != nil {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}

	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, *enc.TodaysID), http.StatusSeeOther)

}

func (su *ImageHandler) UploadUrl(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	var errorHandler ErrorHandler

	sm := api.GetBaseWeb(r)
	category := r.FormValue("category")
	url := r.FormValue("urlTwitch")
	title := r.FormValue("title")
	description := r.FormValue("description")

	if url == "" {
		errorHandler.PageErrorMSG(http.StatusNotAcceptable, configuration.ErrorEmptyField, configuration.ROUTER_UPLOAD_IMAGES, w, r, sm)
		return
	}
	var extencionURL = map[string]bool{
		"Twitch":  true,
		"Youtube": true,
	}
	_, ok := extencionURL[category]
	if !ok {
		errorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, sm)
		return
	}
	/*
		example
		https://clips.twitch.tv/BetterSarcasticZebraFreakinStinkin-rhi_YIaUefuRKPwb
	*/

	var accController controller.AccountController
	acc, err := accController.GetAccount(sm.Email)
	if err != nil {
		errorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, configuration.ROUTER_INDEX, w, r, sm)
		return
	}
	result := strings.Split(url, "/")
	clipLink := result[len(result)-1] //Identifico el ultimo enlace
	var todaysController controller.TodaysController
	todays, err := todaysController.CreateTodays(models.Todays{
		Title:       title,
		Description: description,
		Account:     *acc,
		AccountID:   acc.ID,
	})
	if err != nil {
		errorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, sm)
		return
	}
	var fileController controller.FileController
	_, err = fileController.CreateFile(models.Files{
		PathConsume: fmt.Sprintf(configuration.TWITCH_CLIPS, clipLink, configuration.Config.Server.Ip),
		PathEncrypt: api.GenerateHash(fmt.Sprintf("%d", todays.ID)),
		TypeFile:    category,
		Todays:      &todays,
		TodaysID:    &todays.ID,
	})
	if err != nil {
		log.Println(err)
		errorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, sm)
		return
	}
	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, todays.ID), http.StatusSeeOther)
}

func (su *ImageHandler) GetPhotosHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	var todaysManager controller.TodaysController
	photos, err := todaysManager.GetTodayPage(page)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(photos)
}
