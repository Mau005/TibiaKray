package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ToolsController struct{}

type SharedLoot struct {
	DateStart time.Time
	DateEnd   time.Time
	Time      time.Time
	LootType  string
	Loot      int
	Supplies  int
	Balance   int
	Character []Character
}

type Character struct {
	Name     string
	Loot     int
	Supplies int
	Balance  int
	Damage   int
	Healing  int
}

func (tc *ToolsController) SharedLoot(content string) error {
	shared_loot, _ := tc.PreparingSharedLoot(content)
	//loot := shared_loot.Loot

	profitXPersona := shared_loot.Balance / len(shared_loot.Character)
	var lootLow []Character
	var lootHight []Character
	for _, elements := range shared_loot.Character {
		elements.Balance = (elements.Supplies + profitXPersona) - elements.Loot
		if 0 >= elements.Balance {
			lootHight = append(lootHight, elements)
		} else {
			lootLow = append(lootLow, elements)
		}

	}
	sub_total := 0
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

			fmt.Printf("%s tiene que pagarle a %s un total de: %d\n", characterHight.Name, characterLow.Name, pagar)

		}

	}

	fmt.Println("Balance: ", shared_loot.Balance)
	fmt.Println("Profit: ", profitXPersona)
	fmt.Println("Total general: ", sub_total)

	return nil
}

func (tc *ToolsController) PreparingSharedLoot(content string) (SharedLoot, error) {

	list_conte := strings.Split(content, "\n")
	var sl SharedLoot
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
	var character_list []Character
	const (
		Name = iota
		Loot
		Supplies
		Balance
		Damage
		Healing
	)
	for i := 0; i < len(character_content["characters"]); i += 6 {
		if len(character_content["characters"]) <= Healing+i {
			break
		}
		name := tc.NormalizeString(character_content["characters"][Name+i], normalize)
		leader := " (Leader)"
		if len(name) >= len(leader){
		  captured := name[len(name)-len(leader):]
		  fmt.Printf("Captured: |%s|/n",captured)
		  if captured == leader{
		    name = strings.Trim(name[:len(name)-len(leader)], " ")
		  }
		}
		loot := tc.normalizeStrInt(1, character_content["characters"][Loot+i], ":", normalize)
		supplies := tc.normalizeStrInt(1, character_content["characters"][Supplies+i], ":", normalize)
		balance := tc.normalizeStrInt(1, character_content["characters"][Balance+i], ":", normalize)
		damage := tc.normalizeStrInt(1, character_content["characters"][Damage+i], ":", normalize)
		healing := tc.normalizeStrInt(1, character_content["characters"][Healing+i], ":", normalize)
		character_list = append(character_list, Character{
			Name:     name,
			Loot:     loot,
			Supplies: supplies,
			Balance:  balance,
			Damage:   damage,
			Healing:  healing})
	}
	sl.Character = character_list

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

func (tc *ToolsController) SharedExp(level string) (string, error){
  levelBase, err := strconv.ParseInt(level, 10, 32)
  if err != nil || levelBase == 0 {
    return "", err
  }
  porcent := levelBase / 3
  min := levelBase - porcent
  max := (min * 1.5) + levelBase
  return fmt.Sprintf("Min: %d, Max: %d", min, max), nil
}