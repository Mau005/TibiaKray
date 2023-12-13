package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	conf "github.com/Mau005/KraynoSerer/configuration"
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
			log.Println("entro?")
			go creaturesController.CollectorCreature()
		case "bosses":
			go creaturesController.CollectorBosses()
		case "allcreatures":
			go creaturesController.CollectorCreature()
			go creaturesController.CollectorBosses()
		}
	}
	/*

		creaturesController.CollectorCreature()
		creaturesController.CollectorBosses()
	*/
	log.Println("Listening Server Run", fmt.Sprintf("%s:%d", conf.Config.Server.Ip, conf.Config.Server.Port))
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", conf.Config.Server.Ip, conf.Config.Server.Port),
		router.NewRouter()))

}
