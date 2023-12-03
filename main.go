package main

import (
	"fmt"
	"log"
	"net/http"

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

	log.Println("Listening Server Run", fmt.Sprintf("%s:%d", conf.Config.Server.Ip, conf.Config.Server.Port))
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", conf.Config.Server.Ip, conf.Config.Server.Port),
		router.NewRouter()))

}
