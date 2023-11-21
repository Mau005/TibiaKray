package handler

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
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

	acc, err := api.GetSessionAccount(r)
	if err != nil {
		//Problem Security solicitude log in web
		log.Println("Usuario no se encuentra logueado no esta autorizado")
		return
	}

	var extencion = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"bmp":  true,
	}
	title := r.FormValue("title")
	description := r.FormValue("description")
	file, handler, err := r.FormFile("documents")

	if err != nil {
		log.Println("No se pudo cargar el file")
		return
	}
	defer file.Close()

	verify := strings.Split(handler.Filename, ".")

	if !(extencion[verify[len(verify)-1]]) {
		log.Println("La extencion indicada no es correcta")
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
		log.Println("No se creo el todays correctamente")
		return
	}

	newName := api.GenerateHash(fmt.Sprintf("%d", todays.ID)) //Genero un sha256, proceso el id para ese hash, luego sumo la extencon del file
	fileExtencion := newName + "." + verify[len(verify)-1]
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
		//En esta parte debo eliminar el todays generado
		return
	}

	var encr models.Files
	encr.TodaysID = &todays.ID
	encr.PathConsume = pathConsume
	encr.PathEncrypt = fileExtencion
	encr.PathOrigin = pathSave + fileExtencion
	encr.Todays = &todays

	var encrController controller.FileEncryptsController

	enc, err := encrController.CreateEncrypFile(encr)
	if err != nil {
		return
	}

	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, *enc.TodaysID), http.StatusSeeOther)

}
