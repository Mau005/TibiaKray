package controller

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Mau005/KraynoSerer/models"
)

type ManagerController struct {
	PendingsPlayers map[string]models.Player
}

func NewManagerController() (*ManagerController, error) {
	var ManagerController ManagerController
	ManagerController.PendingsPlayers = make(map[string]models.Player)

	return &ManagerController, nil
}
func (mc *ManagerController) GetPlayer(name string) (models.Player, error) {

	player, ok := mc.PendingsPlayers[name]
	if !ok {
		return models.Player{}, errors.New("Personaje no existe en la tabla de activación")
	}
	return player, nil
}

func (mc *ManagerController) Update() {
	var collector CollectorController
	var playerManager PlayerController
	for {
		var keysToDelete []string
		for namePlayerKey, playerTemp := range mc.PendingsPlayers {
			if playerTemp.CountTimeActive != nil && playerTemp.CountTimeActive.After(time.Now()) {
				// Agrega la clave a la lista de claves a eliminar
				playerWeb, err := collector.GetPlayer(namePlayerKey)
				if err != nil {
					log.Println("Problemas de protocolo ", err.Error())
					continue
				}
				regexBase := `\|TibiaKray=([^|]+)\|`
				uuidWebPLayer, err := mc.SearchRegex(regexBase, playerWeb.Comment)
				if err != nil {
					//log.Println(err) indica que no encontro nad.
					continue
				}

				if playerTemp.Uuid == uuidWebPLayer {
					player, err := playerManager.GetPlayer(playerTemp.Name)
					if err != nil {
						player, err = playerManager.CreatePlayer(playerTemp)
						if err != nil {
							log.Println(err)
						}
					} else {
						player.AccountID = playerTemp.AccountID
						player, err = playerManager.UpdatePLayer(player)
						if err != nil {
							log.Println(err)
						}
					}
					log.Println("Se proceso el personaje y se eliminara del temporal: ", namePlayerKey)
					keysToDelete = append(keysToDelete, player.Name)

				}
			} else {
				fmt.Println("Ya se perdio el tiempo establecido: ", namePlayerKey)
				keysToDelete = append(keysToDelete, namePlayerKey)
			}

		}
		for _, keys := range keysToDelete {
			delete(mc.PendingsPlayers, keys)
		}
		time.Sleep(1 * time.Minute)
	}
}

func (mc *ManagerController) SearchRegex(exprecion, content string) (string, error) {
	re := regexp.MustCompile(exprecion)
	result := re.FindStringSubmatch(content)
	if len(result) > 1 {
		return strings.Trim(result[1], " "), nil
	}

	return "", errors.New("Error de procesado, no encontro nada")
}

func (mc *ManagerController) AddPlayer(player models.Player) models.Player {
	var api ApiController
	player.Uuid = api.GenerateUUid()
	fmt.Printf("Se agrego: %s con key: %s\n", player.Name, player.Uuid)

	timeNow := time.Now().Add(time.Minute * 30)
	player.CountTimeActive = &timeNow
	mc.PendingsPlayers[player.Name] = player
	return player
}
