package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/router"
)

func main() {
	var api controller.ApiController
	err := api.InitServices()
	if err != nil {
		log.Panic(err)
	}

	if len(os.Args) > 1 {
		var creaturesController controller.EntitysCreatures
		command := os.Args[1]
		switch command {
		case "creature":
			go creaturesController.CollectorCreature()
		case "bosses":
			go creaturesController.CollectorBosses()
		case "allcreatures":
			go creaturesController.CollectorCreature()
			go creaturesController.CollectorBosses()
		}
	}

	server := &http.Server{
		Addr:    ":8000",
		Handler: router.NewRouter(),
	}

	certFile := "/etc/letsencrypt/live/tibiakray.info/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/tibiakray.info/privkey.pem"

	log.Println("Iniciando el servidor HTTPS en el puerto 8000")
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatal("Error al iniciar el servidor TLS: ", err)
	}
}
