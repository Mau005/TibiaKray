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

func (su *ImageHandler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("static/upload_files.html")
	if err != nil {
		return
	}
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	templ.Execute(w, sm)
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

	var encrController controller.FileEncryptsController

	enc, err := encrController.CreateEncrypFile(encr)
	if err != nil {
		base := api.GetBaseWeb(r)
		log.Println(err)
		ErrorHandler.PageErrorMSG(http.StatusConflict, configuration.ErrorInternal, configuration.ROUTER_UPLOAD_IMAGES, w, r, base)
		return
	}

	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, *enc.TodaysID), http.StatusSeeOther)

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
	json.NewEncoder(w).Encode(photos)
}
