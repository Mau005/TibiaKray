package controller

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/models"
)

type ToolsController struct{}

func (tc *ToolsController) SharedLoot(content string) (models.SharedLoot, map[string][]string, error) {
	shared_loot, _ := tc.PreparingSharedLoot(content)
	//loot := shared_loot.Loot
	var paymentsNames map[string][]string
	if shared_loot.Balance == 0 || len(shared_loot.Character) <= 1 {
		return shared_loot, paymentsNames, errors.New("Error de protoclo")
	}

	profitXPersona := shared_loot.Balance / len(shared_loot.Character)
	var lootLow []models.CharacterShared
	var lootHight []models.CharacterShared
	for _, elements := range shared_loot.Character {
		elements.Balance = (elements.Supplies + profitXPersona) - elements.Loot
		if 0 >= elements.Balance {
			lootHight = append(lootHight, elements)
		} else {
			lootLow = append(lootLow, elements)
		}

	}
	paymentsNames = make(map[string][]string, len(shared_loot.Character))
	for indexHight, characterHight := range lootHight {
		for indexLow, characterLow := range lootLow {
			if characterHight.Balance == 0 || characterLow.Balance == 0 {
				continue
			}
			//fmt.Printf("Hight: %d Low:%d\n", characterHight.Balance, characterLow.Balance)
			residuo := characterHight.Balance + characterLow.Balance
			pagar := characterLow.Balance
			//fmt.Println("Residuo ", residuo)
			if residuo >= 0 {
				pagar = -characterHight.Balance
				characterLow.Balance = residuo
				characterHight.Balance = 0

			} else {
				characterLow.Balance = 0
				characterHight.Balance = residuo
			}
			lootHight[indexHight] = characterHight
			lootLow[indexLow] = characterLow
			paymentsNames[characterHight.Name] = append(paymentsNames[characterHight.Name], fmt.Sprintf("transfer %d to %s", pagar, characterLow.Name))
			//"transfer 3434 to Mau Tank Persistente"
			//fmt.Printf("%s tiene que pagarle a %s un total de: %d\n", characterHight.Name, characterLow.Name, pagar)
		}

	}

	return shared_loot, paymentsNames, nil
}

func (tc *ToolsController) PreparingSharedLoot(content string) (models.SharedLoot, error) {

	list_conte := strings.Split(content, "\n")
	var sl models.SharedLoot
	map_contet := make(map[string]string, len(list_conte)) //Map content var init
	normalize := []string{"\r", "\t"}                      //list normalize string
	character_content := make(map[string][]string)         //list do character content
	for index, value := range list_conte {
		//fmt.Printf("I: %d V: %q\n", index, value)
		attribute := strings.Split(value, ":")

		if !(len(attribute) > 1) {
			//Iteramos el siguiente si este no es soportado
			continue
		}
		key := tc.NormalizeString(attribute[0], normalize)      //normalizamos las key
		keyValue := tc.NormalizeString(attribute[1], normalize) //normalizamos los valores
		map_contet[key] = keyValue

		if key == "Balance" {
			//fmt.Println("En el index: ", index, " Se encuentra el Balance")
			character_content["characters"] = list_conte[index+1:] //corto el index correpsondiente a la llegada de character
			break
		}
	}
	//ASing SharedLoot
	sl.Balance = tc.normalizeStrInt(0, map_contet["Balance"], ":", normalize)
	sl.Loot = tc.normalizeStrInt(0, map_contet["Loot"], ":", normalize)
	sl.LootType = map_contet["Loot Type"]
	sl.Supplies = tc.normalizeStrInt(0, map_contet["Supplies"], ":", normalize)
	//END Asing Shared Loot
	var character_list []models.CharacterShared
	const (
		Name = iota
		Loot
		Supplies
		Balance
		Damage
		Healing
	)
	for i := 0; i < len(character_content["characters"]); i += 6 {
		var leaderStatus bool
		if len(character_content["characters"]) <= Healing+i {
			break
		}
		name := tc.NormalizeString(character_content["characters"][Name+i], normalize)
		leader := " (Leader)"
		if len(name) >= len(leader) {
			captured := name[len(name)-len(leader):]
			if captured == leader {
				name = strings.Trim(name[:len(name)-len(leader)], " ")
				leaderStatus = true
			}
		}
		loot := tc.normalizeStrInt(1, character_content["characters"][Loot+i], ":", normalize)
		supplies := tc.normalizeStrInt(1, character_content["characters"][Supplies+i], ":", normalize)
		balance := tc.normalizeStrInt(1, character_content["characters"][Balance+i], ":", normalize)
		damage := tc.normalizeStrInt(1, character_content["characters"][Damage+i], ":", normalize)
		healing := tc.normalizeStrInt(1, character_content["characters"][Healing+i], ":", normalize)
		objeCharacter := models.CharacterShared{
			Name:     name,
			Loot:     loot,
			Supplies: supplies,
			Balance:  balance,
			Damage:   damage,
			Healing:  healing}
		character_list = append(character_list, objeCharacter)
		if leaderStatus {
			sl.Leader = objeCharacter
			leaderStatus = false
		}
	}
	sl.Character = character_list
	if configuration.SharedLootHightNow.Balance <= sl.Balance {
		configuration.SharedLootHightNow = sl
	}

	return sl, nil
}

func (tc *ToolsController) normalizeStrInt(priority int, content, keySplit string, removeList []string) int {
	content = strings.ReplaceAll(content, ",", "")
	containerListProcesing := strings.Split(content, keySplit)
	result := ""

	if len(containerListProcesing) == 1 {
		result = strings.Trim(tc.NormalizeString(containerListProcesing[0], removeList), " ")
		procesingInt, _ := strconv.ParseInt(result, 10, 64)

		return int(procesingInt)
	}
	if priority >= len(containerListProcesing) {
		priority = len(containerListProcesing) - 1
	}
	result = strings.Trim(tc.NormalizeString(containerListProcesing[priority], removeList), " ")
	procesingInt, _ := strconv.ParseInt(result, 10, 64)

	return int(procesingInt)
}

func (tc *ToolsController) NormalizeString(content string, removeList []string) string {
	for _, valueContent := range removeList {
		content = strings.ReplaceAll(content, valueContent, "")
	}
	return strings.Trim(content, " ")
}

func (tc *ToolsController) SharedExp(level string, sm models.StructModel) (string, error) {
	levelBase, err := strconv.ParseInt(level, 10, 32)
	if err != nil || levelBase == 0 {
		return "", err
	}
	max := int(math.Round(float64(levelBase) * 1.5))
	min := int(math.Round(float64(levelBase) / 1.5))

	return fmt.Sprintf("%s: %d, %s: %d", Lenguaje[sm.LenguajeDefault]["min"], min, Lenguaje[sm.LenguajeDefault]["max"], max), nil
}

func (tc *ToolsController) InitRashid() string {

	date := time.Now()
	//1 es lunes tiene que ser lunes
	day := int(date.Weekday())
	if !(configuration.Config.Server.ServerSave <= date.Hour()) {
		fmt.Println("Entro?")
		day -= 1
	}
	if day == -1 {
		day = 6
	}

	switch day {
	case configuration.Sunday:
		configuration.Rashid = "RashidSunday"
	case configuration.Monday:
		configuration.Rashid = "RashidMonday"
	case configuration.Tuesday:
		configuration.Rashid = "RashidTuesday"
	case configuration.Wednesday:
		configuration.Rashid = "RashidWednesday"
	case configuration.Thursday:
		configuration.Rashid = "RashidThursday"
	case configuration.Friday:
		configuration.Rashid = "RashidFriday"
	case configuration.Saturday:
		configuration.Rashid = "RashidSaturday"

	}
	fmt.Println(configuration.Friday)
	return configuration.Rashid

}
