package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gorilla/mux"
)

type CreaturesHandler struct{}

func (ch *CreaturesHandler) CreaturesHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var creaturesController controller.EntitysCreatures
	creatures, err := creaturesController.GetCreatures()

	structNew := struct {
		models.StructModel
		Creatures []models.Creatures
	}{
		StructModel: sm,
		Creatures:   creatures,
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_CREATURES)
	if err != nil {
		log.Println(err)
		return
	}

	templ.Execute(w, structNew)
}

func (ch *CreaturesHandler) BossesHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	var creaturesController controller.EntitysCreatures
	Bosses, err := creaturesController.GetBosses()
	structNew := struct {
		models.StructModel
		Bosses []models.Bosses
	}{
		StructModel: sm,
		Bosses:      Bosses,
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_BOSSES)
	if err != nil {
		log.Println(err)
		return
	}

	templ.Execute(w, structNew)
}

func (ch *CreaturesHandler) CreaturesIdHandler(w http.ResponseWriter, r *http.Request) {

	args := mux.Vars(r)

	idCreatures, err := strconv.Atoi(args["id"])
	if err != nil {
		log.Println(err)
		return
	}
	var creaturesController controller.EntitysCreatures
	creature, err := creaturesController.GetIdCreatures(uint(idCreatures))
	if err != nil {
		log.Println(err)
		return
	}

	creatures, err := creaturesController.GetThreeCreaturesRand()
	if err != nil {
		log.Println(err)
		return
	}

	var api controller.ApiController
	sm := api.GetBaseWeb(r)
	NameDamage, MaxDamage := creaturesController.MaxDamageChanges(creature)
	structNew := struct {
		models.StructModel
		Creature   models.Creatures
		Creatures  []models.Creatures
		MaxDamage  uint
		NameDamage string
		TypeWeb    string
	}{
		StructModel: sm,
		Creature:    creature,
		Creatures:   creatures,
		MaxDamage:   MaxDamage,
		NameDamage:  NameDamage,
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_CREATURES_ID)
	if err != nil {
		log.Println(err)
		return
	}

	templ.Execute(w, structNew)
}

func (ch *CreaturesHandler) CreaturesPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id_creature")
	idCreature, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid creature ID", http.StatusBadRequest)
		return
	}

	var creatureController controller.EntitysCreatures
	creature, err := creatureController.GetIdCreatures(uint(idCreature))
	if err != nil {
		log.Println(err)
		http.Error(w, "Creature not found", http.StatusNotFound)
		return
	}

	// Parse form values
	healthStr := r.FormValue("health")
	experienceStr := r.FormValue("experience")
	armorStr := r.FormValue("armor")
	hasteStr := r.FormValue("haste")
	charmStr := r.FormValue("charm")
	mitigationStr := r.FormValue("mitigation")
	dificultyStr := r.FormValue("dificulty")
	locationStr := r.FormValue("location")
	race := r.FormValue("race")
	physicalStr := r.FormValue("physical")
	earthStr := r.FormValue("earth")
	fireStr := r.FormValue("fire")
	deathStr := r.FormValue("death")
	energyStr := r.FormValue("energy")
	holyStr := r.FormValue("holy")
	iceStr := r.FormValue("ice")
	healingStr := r.FormValue("healing")
	maxdamageStr := r.FormValue("maxdamage")
	maxdamage, err := strconv.ParseUint(maxdamageStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid MaxDamge value", http.StatusBadRequest)
	}
	creature.MaxDamage = uint(maxdamage)

	// Parse numerical values
	physical, err := strconv.ParseUint(physicalStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid physical value", http.StatusBadRequest)
		return
	}

	earth, err := strconv.ParseUint(earthStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid earth value", http.StatusBadRequest)
		return
	}

	fire, err := strconv.ParseUint(fireStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid fire value", http.StatusBadRequest)
		return
	}

	death, err := strconv.ParseUint(deathStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid death value", http.StatusBadRequest)
		return
	}

	energy, err := strconv.ParseUint(energyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid energy value", http.StatusBadRequest)
		return
	}

	holy, err := strconv.ParseUint(holyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid holy value", http.StatusBadRequest)
		return
	}

	ice, err := strconv.ParseUint(iceStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ice value", http.StatusBadRequest)
		return
	}

	healing, err := strconv.ParseUint(healingStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid healing value", http.StatusBadRequest)
		return
	}

	// Parse numerical values
	health, err := strconv.ParseUint(healthStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid health value", http.StatusBadRequest)
		return
	}

	experience, err := strconv.ParseUint(experienceStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid experience value", http.StatusBadRequest)
		return
	}

	armor, err := strconv.ParseUint(armorStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid armor value", http.StatusBadRequest)
		return
	}

	haste, err := strconv.ParseUint(hasteStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid haste value", http.StatusBadRequest)
		return
	}

	charm, err := strconv.ParseUint(charmStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid charm value", http.StatusBadRequest)
		return
	}

	mitigation, err := strconv.ParseFloat(mitigationStr, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid mitigation value", http.StatusBadRequest)
		return
	}

	dificulty, err := strconv.ParseUint(dificultyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid difficulty value", http.StatusBadRequest)
		return
	}

	// Update creature fields
	creature.Health = uint(health)
	creature.Experience = uint(experience)
	creature.Armor = uint(armor)
	creature.Haste = uint(haste)
	creature.Charm = uint(charm)
	creature.Mitigation = float32(mitigation)
	creature.Dificulty = uint(dificulty)
	creature.Locations = locationStr
	creature.Race = race
	creature.Physical = uint(physical)
	creature.Earth = uint(earth)
	creature.Fire = uint(fire)
	creature.Death = uint(death)
	creature.Energy = uint(energy)
	creature.Holy = uint(holy)
	creature.Ice = uint(ice)
	creature.Healing = uint(healing)
	creature.Loot = r.FormValue("loot")

	// Update boolean fields
	creature.PushObject = r.FormValue("pushobj") == "on"
	creature.SummonConvince = r.FormValue("summonconvince") == "on"
	creature.Pushable = r.FormValue("pushable") == "on"
	creature.Paralyzable = r.FormValue("paralyzable") == "on"

	// Save the updated creature
	creature, err = creatureController.SaveCreature(creature)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error saving creature", http.StatusInternalServerError)
		return
	}

	// Redirect to the creature's page
	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_CREATURE_ID, creature.ID), http.StatusSeeOther)
}

func (ch *CreaturesHandler) BossesIDHandler(w http.ResponseWriter, r *http.Request) {

	args := mux.Vars(r)

	idCreatures, err := strconv.Atoi(args["id"])
	if err != nil {
		log.Println(err)
		return
	}
	var creaturesController controller.EntitysCreatures
	creature, err := creaturesController.GetIdBosses(uint(idCreatures))
	if err != nil {
		log.Println(err)
		return
	}

	creatures, err := creaturesController.GetThreeBossRand()
	if err != nil {
		log.Println(err)
		return
	}

	var api controller.ApiController
	sm := api.GetBaseWeb(r)
	NameDamage, MaxDamage := creaturesController.MaxDamageChangesBoss(creature)
	structNew := struct {
		models.StructModel
		Creature   models.Bosses
		Creatures  []models.Bosses
		MaxDamage  uint
		NameDamage string
		TypeWeb    string
	}{
		StructModel: sm,
		Creature:    creature,
		Creatures:   creatures,
		MaxDamage:   MaxDamage,
		NameDamage:  NameDamage,
		TypeWeb:     "Boss",
	}

	templ, err := template.ParseFiles(configuration.PATH_WEB_CREATURES_ID)
	if err != nil {
		log.Println(err)
		return
	}

	templ.Execute(w, structNew)
}

func (ch *CreaturesHandler) BossPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id_creature")
	idCreature, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid creature ID", http.StatusBadRequest)
		return
	}

	var creatureController controller.EntitysCreatures
	creature, err := creatureController.GetIdBosses(uint(idCreature))
	if err != nil {
		log.Println(err)
		http.Error(w, "Creature not found", http.StatusNotFound)
		return
	}

	// Parse form values
	healthStr := r.FormValue("health")
	experienceStr := r.FormValue("experience")
	armorStr := r.FormValue("armor")
	hasteStr := r.FormValue("haste")
	charmStr := r.FormValue("charm")
	mitigationStr := r.FormValue("mitigation")
	dificultyStr := r.FormValue("dificulty")
	locationStr := r.FormValue("location")
	race := r.FormValue("race")
	physicalStr := r.FormValue("physical")
	earthStr := r.FormValue("earth")
	fireStr := r.FormValue("fire")
	deathStr := r.FormValue("death")
	energyStr := r.FormValue("energy")
	holyStr := r.FormValue("holy")
	iceStr := r.FormValue("ice")
	healingStr := r.FormValue("healing")
	maxdamageStr := r.FormValue("maxdamage")
	maxdamage, err := strconv.ParseUint(maxdamageStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid MaxDamge value", http.StatusBadRequest)
	}
	creature.MaxDamage = uint(maxdamage)

	// Parse numerical values
	physical, err := strconv.ParseUint(physicalStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid physical value", http.StatusBadRequest)
		return
	}

	earth, err := strconv.ParseUint(earthStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid earth value", http.StatusBadRequest)
		return
	}

	fire, err := strconv.ParseUint(fireStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid fire value", http.StatusBadRequest)
		return
	}

	death, err := strconv.ParseUint(deathStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid death value", http.StatusBadRequest)
		return
	}

	energy, err := strconv.ParseUint(energyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid energy value", http.StatusBadRequest)
		return
	}

	holy, err := strconv.ParseUint(holyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid holy value", http.StatusBadRequest)
		return
	}

	ice, err := strconv.ParseUint(iceStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ice value", http.StatusBadRequest)
		return
	}

	healing, err := strconv.ParseUint(healingStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid healing value", http.StatusBadRequest)
		return
	}

	// Parse numerical values
	health, err := strconv.ParseUint(healthStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid health value", http.StatusBadRequest)
		return
	}

	experience, err := strconv.ParseUint(experienceStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid experience value", http.StatusBadRequest)
		return
	}

	armor, err := strconv.ParseUint(armorStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid armor value", http.StatusBadRequest)
		return
	}

	haste, err := strconv.ParseUint(hasteStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid haste value", http.StatusBadRequest)
		return
	}

	charm, err := strconv.ParseUint(charmStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid charm value", http.StatusBadRequest)
		return
	}

	mitigation, err := strconv.ParseFloat(mitigationStr, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid mitigation value", http.StatusBadRequest)
		return
	}

	dificulty, err := strconv.ParseUint(dificultyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid difficulty value", http.StatusBadRequest)
		return
	}

	// Update creature fields
	creature.Health = uint(health)
	creature.Experience = uint(experience)
	creature.Armor = uint(armor)
	creature.Haste = uint(haste)
	creature.Charm = uint(charm)
	creature.Mitigation = float32(mitigation)
	creature.Dificulty = uint(dificulty)
	creature.Locations = locationStr
	creature.Race = race
	creature.Physical = uint(physical)
	creature.Earth = uint(earth)
	creature.Fire = uint(fire)
	creature.Death = uint(death)
	creature.Energy = uint(energy)
	creature.Holy = uint(holy)
	creature.Ice = uint(ice)
	creature.Healing = uint(healing)
	creature.Loot = r.FormValue("loot")

	// Update boolean fields
	creature.PushObject = r.FormValue("pushobj") == "on"
	creature.SummonConvince = r.FormValue("summonconvince") == "on"
	creature.Pushable = r.FormValue("pushable") == "on"
	creature.Paralyzable = r.FormValue("paralyzable") == "on"

	// Save the updated creature
	creature, err = creatureController.SaveBoss(creature)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error saving creature", http.StatusInternalServerError)
		return
	}

	// Redirect to the creature's page
	http.Redirect(w, r, fmt.Sprintf(configuration.ROUTER_BOSSES_ID, creature.ID), http.StatusSeeOther)
}
