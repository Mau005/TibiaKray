package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mau005/KraynoSerer/configuration"
	conf "github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/Mau005/KraynoSerer/router"
)

func main() {

	err := conf.LoadConfiguration("config.yml")
	if err != nil {
		log.Panic(err)
	}

	var cc controller.CollectorController
	err = cc.GenerateNewsTibia()
	if err != nil {
		log.Println(fmt.Sprintf("[NEWS][ERROR] %s", err.Error()))
	}
	log.Println("[NEWS] Load for Tibia.com")

	var toolsManager controller.ToolsController
	toolsManager.InitRashid()

	go func() {
		for {
			ResetDefaultWeb()
		}
	}()
	var api controller.ApiController
	err = api.InitLenguaje("data/lenguaje.csv")
	if err != nil {
		log.Panic(err)
	}

	err = api.InitLenguajeServer("data/errorServer.csv")
	if err != nil {
		log.Panic(err)
	}

	err = database.ConnectionDataBase()
	if err != nil {
		log.Panic(err)
	}
	/*
		if len(os.Args) >1 {
			command := os.Args[1]
			switch command{
			case "generate":
				database.Generate()
			}
		}
	*/

	log.Println("Listening Server Run", fmt.Sprintf("%s:%d", conf.Config.Server.Ip, conf.Config.Server.Port))
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", conf.Config.Server.Ip, conf.Config.Server.Port),
		router.NewRouter()))

}

func ResetDefaultWeb() {

	now := time.Now()
	// Calcula la duración hasta la próxima 6 de la mañana
	nextSixAM := time.Date(now.Year(), now.Month(), now.Day(), configuration.Config.Server.ServerSave, 0, 0, 0, now.Location())
	if now.After(nextSixAM) {
		// Si ya es después de las 6 de la mañana hoy, programa para mañana
		nextSixAM = nextSixAM.Add(24 * time.Hour)
	}

	durationUntilSixAM := nextSixAM.Sub(now)

	// Configura un temporizador para ejecutar la función a las 6 de la mañana
	timer := time.NewTimer(durationUntilSixAM)
	<-timer.C // Espera hasta que el temporizador alcance su límite

	// Ejecuta tu función aquí
	log.Println("Reset Data: ", nextSixAM)
	var cc controller.CollectorController
	err := cc.GenerateNewsTibia()
	if err != nil {
		log.Println(fmt.Sprintf("[NEWS][ERROR] %s", err.Error()))
	}
	log.Println("[NEWS] Load for Tibia.com")

	var toolsManager controller.ToolsController
	toolsManager.InitRashid()

	//reset sharedloot
	configuration.SharedLootHightNow = models.SharedLoot{}

}
