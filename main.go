package main

import (
	"fmt"
	"log"
	"net/http"

	conf "github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/database"
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

	var api controller.ApiController
	err = api.InitLenguaje("data/lenguaje.csv")
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
