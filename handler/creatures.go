package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gorilla/mux"
)

type CreaturesHandler struct{}

func (ch *CreaturesHandler) GetCreaturesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	name := q.Get("name")
	race := q.Get("race")
	sort := q.Get("sort")
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 25
	}
	if pageSize > 5000 {
		pageSize = 5000
	}
	var hpMin *uint
	if v := q.Get("hpMin"); v != "" {
		n, _ := strconv.Atoi(v)
		u := uint(n)
		hpMin = &u
	}

	var hpMax *uint
	if v := q.Get("hpMax"); v != "" {
		n, _ := strconv.Atoi(v)
		u := uint(n)
		hpMax = &u
	}

	var mitMin *float32
	if v := q.Get("mitMin"); v != "" {
		f, _ := strconv.ParseFloat(v, 32)
		m := float32(f)
		mitMin = &m
	}

	parseBool := func(key string) *bool {
		if q.Get(key) == "" {
			return nil
		}
		b := q.Get(key) == "true"
		return &b
	}

	list, total, pageOut, pageSizeOut, err := ch.GetCreaturesFiltered(
		name,
		race,
		hpMin,
		hpMax,
		mitMin,
		parseBool("pushable"),
		parseBool("paralyzable"),
		parseBool("pushObject"),
		parseBool("summonConvince"),
		sort,
		page,
		pageSize,
	)
	pages := int64((total + int64(pageSizeOut) - 1) / int64(pageSizeOut))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(pageOut, pageSizeOut)

	resp := map[string]any{
		"total":    total,
		"page":     pageOut,
		"pageSize": pageSizeOut,
		"pages":    pages,
		"data":     list,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (ch *CreaturesHandler) GetCreaturesFiltered(
	name string,
	race string,
	hpMin *uint,
	hpMax *uint,
	mitMin *float32,
	pushable *bool,
	paralyzable *bool,
	pushObject *bool,
	summonConvince *bool,
	sort string,
	page int,
	pageSize int,
) (list []models.Creatures, total int64, pageOut int, pageSizeOut int, err error) {

	db := database.DB.Model(&models.Creatures{})

	// ---- Filters ----
	if name != "" {
		name = strings.ToLower(strings.TrimSpace(name))
		db = db.Where("LOWER(name) LIKE ?", "%"+name+"%")
	}
	if race != "" {
		db = db.Where("race = ?", race)
	}
	if hpMin != nil {
		db = db.Where("health >= ?", *hpMin)
	}
	if hpMax != nil {
		db = db.Where("health <= ?", *hpMax)
	}
	if mitMin != nil {
		db = db.Where("mitigation >= ?", *mitMin)
	}
	if pushable != nil {
		db = db.Where("pushable = ?", *pushable)
	}
	if paralyzable != nil {
		db = db.Where("paralyzable = ?", *paralyzable)
	}
	if pushObject != nil {
		db = db.Where("push_object = ?", *pushObject)
	}
	if summonConvince != nil {
		db = db.Where("summon_convince = ?", *summonConvince)
	}

	// ---- Count total ----
	if err = db.Count(&total).Error; err != nil {
		return
	}

	// ---- Sorting ----
	switch sort {
	case "name.asc":
		db = db.Order("name ASC")
	case "name.desc":
		db = db.Order("name DESC")
	case "health.asc":
		db = db.Order("health ASC")
	case "health.desc":
		db = db.Order("health DESC")
	case "exp.asc":
		db = db.Order("experience ASC")
	case "exp.desc":
		db = db.Order("experience DESC")
	case "armor.asc":
		db = db.Order("armor ASC")
	case "armor.desc":
		db = db.Order("armor DESC")
	default:
		db = db.Order("name ASC")
	}

	// ---- Pagination (validación + outputs) ----
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 25
	}
	if pageSize > 5000 {
		pageSize = 5000
	}

	pageOut = page
	pageSizeOut = pageSize

	offset := (page - 1) * pageSize

	err = db.
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error

	return
}

func (ch *CreaturesHandler) CreaturesHandler(w http.ResponseWriter, r *http.Request) {
	var api controller.ApiController
	sm := api.GetBaseWeb(r)

	structNew := struct {
		models.StructModel
	}{
		StructModel: sm,
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
	if err != nil {
		log.Println(err)
	}
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
		MaxDamage  int
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
	physical, err := strconv.ParseInt(physicalStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid physical value", http.StatusBadRequest)
		return
	}

	earth, err := strconv.ParseInt(earthStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid earth value", http.StatusBadRequest)
		return
	}

	fire, err := strconv.ParseInt(fireStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid fire value", http.StatusBadRequest)
		return
	}

	death, err := strconv.ParseInt(deathStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid death value", http.StatusBadRequest)
		return
	}

	energy, err := strconv.ParseInt(energyStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid energy value", http.StatusBadRequest)
		return
	}

	holy, err := strconv.ParseInt(holyStr, 10, 64)
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
	creature.Physical = int(physical)
	creature.Earth = int(earth)
	creature.Fire = int(fire)
	creature.Death = int(death)
	creature.Energy = int(energy)
	creature.Holy = int(holy)
	creature.Ice = int(ice)
	creature.Healing = int(healing)
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
		MaxDamage  int
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
	creature.Physical = int(physical)
	creature.Earth = int(earth)
	creature.Fire = int(fire)
	creature.Death = int(death)
	creature.Energy = int(energy)
	creature.Holy = int(holy)
	creature.Ice = int(ice)
	creature.Healing = int(healing)
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
