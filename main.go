package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mau005/KraynoSerer/configuration"
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
		case "loadLua":
			creaturesController.LoadLuaMonster()
		case "creature":
			go func() {
				creaturesController.CollectorCreature()
				creaturesController.LoadLuaMonster()
			}()

		case "bosses":
			go creaturesController.CollectorBosses()
			creaturesController.LoadLuaMonster()
		case "allcreatures":

			go func() {
				creaturesController.CollectorCreature()
				creaturesController.CollectorBosses()
				creaturesController.LoadLuaMonster()
			}()

		}
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", configuration.Config.Server.Ip, configuration.Config.Server.Port),
		Handler: router.NewRouter(),
	}
	log.Println("HTTP :80 (ACME challenge)")
	if err := server.ListenAndServe(); err != nil {
		log.Println("HTTP server stopped:", err)
	}
	// certFile := "/etc/letsencrypt/live/tibiakray.info/fullchain.pem"
	// keyFile := "/etc/letsencrypt/live/tibiakray.info/privkey.pem"

	// log.Println("Iniciando el servidor HTTPS en el puerto 8000")
	// if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
	// 	log.Fatal("Error al iniciar el servidor TLS: ", err)
	// }
}
