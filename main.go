package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/router"
	"golang.org/x/crypto/acme/autocert"
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

	certCache := autocert.DirCache("certs")

	m := &autocert.Manager{
		Cache:  certCache,
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"tibiakray.info", // replace with your domain
		),
	}

	server := &http.Server{
		Addr:      ":8000", // Cambiado a puerto 8000
		Handler:   router.NewRouter(),
		TLSConfig: m.TLSConfig(),
	}

	log.Println("Starting HTTPS server on port 443...")
	log.Fatal(server.ListenAndServeTLS("", ""))

}
