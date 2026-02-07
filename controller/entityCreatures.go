package controller

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/gocolly/colly/v2"
	lua "github.com/yuin/gopher-lua"
)

type EntitysCreatures struct{}
type ImportResult struct {
	Category string
	File     string
	Creature []byte
}

func categoryFromPath(root, full string) string {
	rel, _ := filepath.Rel(root, full) // Category/monster.lua
	parts := strings.Split(rel, string(os.PathSeparator))
	if len(parts) >= 2 {
		return parts[0]
	}
	return ""
}
func shouldSkip(path string) bool {
	p := strings.ToLower(filepath.ToSlash(path))
	if strings.Contains(p, "/quests/") {
		return true
	}
	if strings.HasSuffix(p, "_functions.lua") { // helpers, no monsters
		return true
	}
	return false
}

func parseMonsterLua(file string) (models.Creatures, error) {
	L := lua.NewState()
	defer L.Close()

	var captured *lua.LTable
	var capturedName string

	// Stubs básicos (para que el script no reviente)
	// 1) Constantes COMBAT_* (si no existen, el script fallará al construir elements)
	defineCombatConstants(L)

	// 2) BESTY_RACE_AQUATIC etc (si aparecen)
	defineBestiaryConstants(L)

	// 3) Game.createMonsterType
	game := L.NewTable()
	L.SetGlobal("Game", game)

	L.SetGlobal("RegisterPrimalPackBeast", L.NewFunction(func(L *lua.LState) int {
		return 0
	}))
	L.SetField(game, "createMonsterType", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		capturedName = name

		mType := L.NewTable()

		// mType:register(monster)
		L.SetField(mType, "register", L.NewFunction(func(L *lua.LState) int {
			monsterTbl := L.CheckTable(2) // self=1, monster=2
			captured = monsterTbl
			return 0
		}))

		L.Push(mType)
		return 1
	}))
	if shouldSkip(file) {
		return models.Creatures{}, nil
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return models.Creatures{}, err
	}

	// 1️⃣ sanitiza strings problemáticos (Locations)
	fixed := sanitizeLocationsToLongString(string(content))

	// Ejecuta el archivo
	if err := L.DoString(fixed); err != nil {
		return models.Creatures{}, err
	}
	if captured == nil {
		return models.Creatures{}, fmt.Errorf("no se capturó register(monster)")
	}

	row := models.Creatures{}
	row.Name = capturedName

	// Campos simples
	row.Experience = uint(getIntField(captured, "experience"))
	row.Health = uint(getIntField(captured, "health"))
	row.Race = getStringField(captured, "race")
	row.Mitigation = float32(getNumberFromNested(captured, "defenses", "mitigation"))
	row.Attack = getStringAttacks(captured, "attacks")
	row.Armor = uint(getIntFromNested(captured, "defenses", "armor"))

	// Outfit/lookType -> lo puedes guardar como string "lookType:451" o resolver a ruta real después
	lookType := getIntFromNested(captured, "outfit", "lookType")
	if lookType > 0 {
		row.ImagePath = fmt.Sprintf("lookType:%d", lookType)
	}

	// Bestiary locations
	row.Locations = getStringFromNested(captured, "Bestiary", "Locations")

	// Flags
	row.Pushable = getBoolFromNested(captured, "flags", "pushable")

	// Summon/Convince: en tu modelo está junto; aquí puedes mapearlo como OR
	summonable := getBoolFromNested(captured, "flags", "summonable")
	convinceable := getBoolFromNested(captured, "flags", "convinceable")
	row.SummonConvince = summonable || convinceable

	// PushObject: en Lua hay canPushItems/canPushCreatures
	canPushItems := getBoolFromNested(captured, "flags", "canPushItems")
	canPushCreatures := getBoolFromNested(captured, "flags", "canPushCreatures")
	row.PushObject = canPushItems || canPushCreatures

	// Paralyzable: en Lua immunities {type="paralyze", condition=true} significa inmune.
	// Tu campo "Paralyzable" (se puede paralizar) => invertido.
	paralyzeImmune := getImmunityCondition(captured, "paralyze")
	if paralyzeImmune {
		row.Paralyzable = false
	} else {
		row.Paralyzable = true
	}

	// Loot -> JSON recomendado
	row.Loot = getStringLoot(getTableField(captured, "loot"))

	// Elements -> a tus columnas int
	applyElements(&row, getTableField(captured, "elements"))

	return row, nil
}

func defineCombatConstants(L *lua.LState) {
	// Les damos números cualquiera (solo se usan como “keys” para comparar)
	constants := map[string]int{
		"COMBAT_PHYSICALDAMAGE": 1,
		"COMBAT_ENERGYDAMAGE":   2,
		"COMBAT_EARTHDAMAGE":    3,
		"COMBAT_FIREDAMAGE":     4,
		"COMBAT_LIFEDRAIN":      5,
		"COMBAT_MANADRAIN":      6,
		"COMBAT_DROWNDAMAGE":    7,
		"COMBAT_ICEDAMAGE":      8,
		"COMBAT_HOLYDAMAGE":     9,
		"COMBAT_DEATHDAMAGE":    10,
		"COMBAT_HEALING":        11,
	}
	for k, v := range constants {
		L.SetGlobal(k, lua.LNumber(v))

	}
}

func defineBestiaryConstants(L *lua.LState) {
	// Muchos datapacks usan constantes BESTY_RACE_*
	L.SetGlobal("BESTY_RACE_AQUATIC", lua.LNumber(1))
	L.SetGlobal("BESTY_RACE_UNDEAD", lua.LNumber(2))
	// agrega más si te salen errores
}

// --- Helpers de lectura Lua ---

func getTableField(t *lua.LTable, key string) *lua.LTable {
	v := t.RawGetString(key)
	if tbl, ok := v.(*lua.LTable); ok {
		return tbl
	}
	return nil
}

func getStringField(t *lua.LTable, key string) string {
	v := t.RawGetString(key)
	if s, ok := v.(lua.LString); ok {
		return string(s)
	}
	return ""
}

func getIntField(t *lua.LTable, key string) int {
	v := t.RawGetString(key)
	if n, ok := v.(lua.LNumber); ok {
		return int(n)
	}
	return 0
}

func getIntFromNested(root *lua.LTable, tableKey, fieldKey string) int {
	tbl := getTableField(root, tableKey)
	if tbl == nil {
		return 0
	}
	return getIntField(tbl, fieldKey)
}

func getStringFromNested(root *lua.LTable, tableKey, fieldKey string) string {
	tbl := getTableField(root, tableKey)
	if tbl == nil {
		return ""
	}
	return getStringField(tbl, fieldKey)
}

func getBoolFromNested(root *lua.LTable, tableKey, fieldKey string) bool {
	tbl := getTableField(root, tableKey)
	if tbl == nil {
		return false
	}
	v := tbl.RawGetString(fieldKey)
	if b, ok := v.(lua.LBool); ok {
		return bool(b)
	}
	return false
}

func getNumberFromNested(root *lua.LTable, tableKey, fieldKey string) float64 {
	tbl := getTableField(root, tableKey)
	if tbl == nil {
		return 0
	}
	v := tbl.RawGetString(fieldKey)
	if n, ok := v.(lua.LNumber); ok {
		return float64(n)
	}
	return 0
}

func getStringAttacks(root *lua.LTable, key string) string {
	attacks := getTableField(root, key)
	if attacks == nil {
		return "[]"
	}

	type Attack struct {
		Name        string `json:"name,omitempty"`
		Interval    int    `json:"interval,omitempty"`
		Chance      int    `json:"chance,omitempty"`
		Type        int    `json:"type,omitempty"`
		MinDamage   int    `json:"minDamage,omitempty"`
		MaxDamage   int    `json:"maxDamage,omitempty"`
		Radius      int    `json:"radius,omitempty"`
		Range       int    `json:"range,omitempty"`
		Length      int    `json:"length,omitempty"`
		Spread      int    `json:"spread,omitempty"`
		Effect      int    `json:"effect,omitempty"`
		ShootEffect int    `json:"shootEffect,omitempty"`
		Target      *bool  `json:"target,omitempty"` // pointer para distinguir "false" vs "no existe"
	}

	var out []Attack

	attacks.ForEach(func(_ lua.LValue, v lua.LValue) {
		t, ok := v.(*lua.LTable)
		if !ok {
			return
		}

		a := Attack{
			Name:        getStringField(t, "name"),
			Interval:    getIntField(t, "interval"),
			Chance:      getIntField(t, "chance"),
			Type:        getIntField(t, "type"),
			MinDamage:   getIntField(t, "minDamage"),
			MaxDamage:   getIntField(t, "maxDamage"),
			Radius:      getIntField(t, "radius"),
			Range:       getIntField(t, "range"),
			Length:      getIntField(t, "length"),
			Spread:      getIntField(t, "spread"),
			Effect:      getIntField(t, "effect"),
			ShootEffect: getIntField(t, "shootEffect"),
		}

		// target puede venir como boolean o no venir
		if tv := t.RawGetString("target"); tv != lua.LNil {
			if b, ok := tv.(lua.LBool); ok {
				val := bool(b)
				a.Target = &val
			}
		}

		out = append(out, a)
	})
	result := ""
	for _, v := range out {
		result += fmt.Sprintf("%s Min:%d Max:%d\n", v.Name, v.MinDamage, v.MaxDamage)
	}
	return result
}

func getImmunityCondition(monster *lua.LTable, immunType string) bool {
	imms := getTableField(monster, "immunities")
	if imms == nil {
		return false
	}

	var immune bool
	imms.ForEach(func(_ lua.LValue, v lua.LValue) {
		entry, ok := v.(*lua.LTable)
		if !ok {
			return
		}
		t := entry.RawGetString("type")
		if s, ok := t.(lua.LString); ok && string(s) == immunType {
			c := entry.RawGetString("condition")
			if b, ok := c.(lua.LBool); ok {
				immune = bool(b)
			}
		}
	})
	return immune
}

func getStringLoot(loot *lua.LTable) string {
	if loot == nil {
		return "[]"
	}

	type LootItem struct {
		Name     string `json:"name,omitempty"`
		Chance   int    `json:"chance,omitempty"`
		MaxCount int    `json:"maxCount,omitempty"`
		Id       int    `json:"id,omitempty"`
	}

	var items []LootItem

	loot.ForEach(func(_ lua.LValue, v lua.LValue) {
		entry, ok := v.(*lua.LTable)
		if !ok {
			return
		}
		it := LootItem{
			Name:     getStringField(entry, "name"),
			Chance:   getIntField(entry, "chance"),
			MaxCount: getIntField(entry, "maxCount"),
			Id:       getIntField(entry, "id"),
		}
		items = append(items, it)
	})
	if len(items) == 0 {
		return ""
	}

	dataStr := make([]string, 0)

	for _, v := range items {
		dataStr = append(dataStr, strings.TrimRight(fmt.Sprintf("%s %s", v.Name, chance_items(uint(v.Chance))), " "))
	}
	result := strings.Join(dataStr, ", ")
	return result
}

func chance_items(value uint) string {
	if value >= 9000 {
		return ""
	} else if value >= 5000 {
		return "(Semi Rare)"
	} else if value >= 1000 {
		return "(Rare)"
	} else {
		return "(Very Rare)"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func applyElements(row *models.Creatures, elements *lua.LTable) {
	if elements == nil {
		return
	}
	elements.ForEach(func(_ lua.LValue, v lua.LValue) {
		entry, ok := v.(*lua.LTable)
		if !ok {
			return
		}
		typ := getIntField(entry, "type")
		percent := getIntField(entry, "percent")

		switch typ {
		case 1: // PHYSICAL
			row.Physical = 100 - percent
		case 2: // ENERGY
			row.Energy = 100 - percent
		case 3: // EARTH
			row.Earth = 100 - percent
		case 4: // FIRE
			row.Fire = 100 - percent
		case 7: // DROWN -> si no tienes columna, podrías mapear a Physical o ignorar
			// ignorar o agregar columna Water/Drown
		case 8: // ICE
			row.Ice = 100 - percent
		case 9: // HOLY
			row.Holy = 100 - percent
		case 10: // DEATH
			row.Death = 100 - percent
		case 11: // HEALING
			row.Healing = percent
		}
	})
}

var reLocations = regexp.MustCompile(`(?m)^\s*Locations\s*=\s*"([^"]*)"\s*,?\s*$`)

func sanitizeLocationsToLongString(src string) string {
	return reLocations.ReplaceAllStringFunc(src, func(line string) string {
		m := reLocations.FindStringSubmatch(line)
		if len(m) != 2 {
			return line
		}
		val := m[1]
		// Evita cerrar ]] accidental (raro, pero por seguridad)
		val = strings.ReplaceAll(val, "]]", "] ]")
		return strings.Replace(line, `"`+m[1]+`"`, `[[`+val+`]]`, 1)
	})
}

func (ec *EntitysCreatures) LoadLuaMonster() error {
	root := "./data/lua/monster"
	var results []ImportResult

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".lua") {
			return nil
		}

		category := categoryFromPath(root, path)

		row, err := parseMonsterLua(path)

		data, err := json.Marshal(row)
		if err != nil {
			fmt.Printf("ERR %s: %v\n", path, err)
			return nil
		}
		if err != nil {
			fmt.Printf("ERR %s: %v\n", path, err)
			return nil
		}
		results = append(results, ImportResult{
			Category: category,
			File:     path,
			Creature: data,
		})
		return nil
	})
	fmt.Printf("Importados: %d\n", len(results))
	for _, v := range results {
		switch strings.ToLower(v.Category) {
		case "bosses":
			creatureLoad := models.Bosses{}
			json.Unmarshal(v.Creature, &creatureLoad)
			creature, err := ec.GetNameBosses(creatureLoad.Name)
			if err != nil {
				continue
			}
			creatureLoad.ID = creature.ID
			creatureLoad.CreatedAt = creature.CreatedAt
			creatureLoad.ImagePath = creature.ImagePath
			if err := database.DB.Save(creatureLoad).Error; err != nil {
				log.Println("error al crear Boss ", creatureLoad)
			}

		default:
			creatureLoad := models.Creatures{}
			json.Unmarshal(v.Creature, &creatureLoad)
			creature, err := ec.GetNameMonster(creatureLoad.Name)
			if err != nil {
				continue
			}
			creatureLoad.ID = creature.ID
			creatureLoad.CreatedAt = creature.CreatedAt
			creatureLoad.ImagePath = creature.ImagePath
			if err := database.DB.Save(creatureLoad).Error; err != nil {
				log.Println("error al crear Creature", creatureLoad)
			}
		}
	}

	return nil
}
func (ec *EntitysCreatures) GetNameBosses(name string) (creat models.Bosses, err error) {
	name = strings.TrimSpace(name)
	err = database.DB.
		Where("name = ?", name).
		First(&creat).Error
	if err == nil {
		return
	}

	var variants []string
	variants = append(variants, name)

	lower := strings.ToLower(name)
	if !strings.HasSuffix(lower, "s") {
		variants = append(variants, name+"s")
	} else {
		variants = append(variants, strings.TrimSuffix(name, "s"))
	}

	err = database.DB.
		Where("name IN ?", variants).
		First(&creat).Error
	if err == nil {
		return
	}

	err = database.DB.
		Where("LOWER(name) LIKE ?", strings.ToLower(name)+"%").
		Order("LENGTH(name) ASC").
		First(&creat).Error
	if err == nil {
		return
	}
	err = database.DB.
		Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%").
		Order("LENGTH(name) ASC").
		First(&creat).Error

	return
}

func (ec *EntitysCreatures) GetNameMonster(name string) (creat models.Creatures, err error) {
	name = strings.TrimSpace(name)
	err = database.DB.
		Where("name = ?", name).
		First(&creat).Error
	if err == nil {
		return
	}

	var variants []string
	variants = append(variants, name)

	lower := strings.ToLower(name)
	if !strings.HasSuffix(lower, "s") {
		variants = append(variants, name+"s")
	} else {
		variants = append(variants, strings.TrimSuffix(name, "s"))
	}

	err = database.DB.
		Where("name IN ?", variants).
		First(&creat).Error
	if err == nil {
		return
	}

	err = database.DB.
		Where("LOWER(name) LIKE ?", strings.ToLower(name)+"%").
		Order("LENGTH(name) ASC").
		First(&creat).Error
	if err == nil {
		return
	}
	err = database.DB.
		Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%").
		Order("LENGTH(name) ASC").
		First(&creat).Error

	return
}

func (ec *EntitysCreatures) CollectorCreature() error {
	c := colly.NewCollector()
	var monsters []models.Creatures
	c.OnHTML("div[style='width: 100px; height: 110px; margin: 0px; float: left;']", func(e *colly.HTMLElement) {
		// Encuentra el enlace dentro del div
		// Encuentra la etiqueta img dentro del div
		imgSrc := e.ChildAttr("img", "src")
		// Encuentra el texto dentro del div
		name := strings.TrimSpace(e.ChildText("div"))

		monsters = append(monsters, models.Creatures{Name: name, ImagePath: imgSrc})
	})
	err := c.Visit(configuration.TIBIA_MONSTER)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, value := range monsters {
		var api ApiController
		sliceNames := strings.Split(value.ImagePath, "/")
		nameImage := sliceNames[len(sliceNames)-1]
		pathSave := fmt.Sprintf(configuration.PATH_STATIC_CREATURES, nameImage)
		fmt.Println("dowloads", value.ImagePath)
		fmt.Println("Save: ", pathSave)
		api.DownloadImage(value.ImagePath, pathSave)
		value.ImagePath = pathSave
		_, err := ec.CreateMonster(value)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (ec *EntitysCreatures) CollectorBosses() error {
	c := colly.NewCollector()
	var bosses []models.Bosses
	c.OnHTML("div[style='width: 100px; height: 110px; margin: 0px; float: left;']", func(e *colly.HTMLElement) {
		// Encuentra el enlace dentro del div
		// Encuentra la etiqueta img dentro del div
		imgSrc := e.ChildAttr("img", "src")
		// Encuentra el texto dentro del div
		name := strings.TrimSpace(e.ChildText("div"))

		bosses = append(bosses, models.Bosses{Creatures: models.Creatures{Name: strings.Trim(name, " "), ImagePath: imgSrc}})
	})

	err := c.Visit(configuration.TIBIA_BOSSES)
	if err != nil {
		log.Println(err)
	}
	for _, value := range bosses {
		var api ApiController
		sliceNames := strings.Split(value.ImagePath, "/")
		nameImage := sliceNames[len(sliceNames)-1]
		pathSave := fmt.Sprintf(configuration.PATH_STATIC_BOSSES, nameImage)
		err := api.DownloadImage(value.ImagePath, pathSave)
		log.Println(err)
		value.ImagePath = pathSave
		_, err = ec.CreateBosses(value)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (ec *EntitysCreatures) GetCreatures() (creatures []models.Creatures, err error) {
	if err = database.DB.Find(&creatures).Error; err != nil {
		return
	}
	return
}
func (ec *EntitysCreatures) GetIdCreatures(idCreature uint) (creatures models.Creatures, err error) {
	if err = database.DB.Where("id = ?", idCreature).First(&creatures).Error; err != nil {
		return
	}
	return
}

func (ec *EntitysCreatures) GetThreeCreaturesRand() (creatures []models.Creatures, err error) {
	var allCreatures []models.Creatures
	database.DB.Find(&allCreatures)

	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	return
}
func (ec *EntitysCreatures) GetThreeBossRand() (creatures []models.Bosses, err error) {
	var allCreatures []models.Bosses
	database.DB.Find(&allCreatures)

	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	return
}

func (ec *EntitysCreatures) GetBosses() (bosses []models.Bosses, err error) {
	if err = database.DB.Find(&bosses).Error; err != nil {
		return
	}
	return
}
func (ec *EntitysCreatures) GetIdBosses(idCreature uint) (creatures models.Bosses, err error) {
	if err = database.DB.Where("id = ?", idCreature).First(&creatures).Error; err != nil {
		return
	}
	return
}

func (ec *EntitysCreatures) CreateMonster(monster models.Creatures) (models.Creatures, error) {
	if err := database.DB.Create(&monster).Error; err != nil {
		return monster, err
	}
	return monster, nil
}

func (ec *EntitysCreatures) CreateBosses(boss models.Bosses) (models.Bosses, error) {
	if err := database.DB.Create(&boss).Error; err != nil {
		return boss, err
	}
	return boss, nil
}
func (ec *EntitysCreatures) GetThreeBossesRand() (creatures []models.Bosses, err error) {
	var allCreatures []models.Bosses
	database.DB.Find(&allCreatures)

	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	creatures = append(creatures, allCreatures[rand.Intn(len(allCreatures))])
	return
}

func (ec *EntitysCreatures) SaveCreature(monster models.Creatures) (models.Creatures, error) {
	if err := database.DB.Save(&monster).Error; err != nil {
		return monster, err
	}
	return monster, nil
}

func (ec *EntitysCreatures) SaveBoss(monster models.Bosses) (models.Bosses, error) {
	if err := database.DB.Save(&monster).Error; err != nil {
		return monster, err
	}
	return monster, nil
}
func (ec *EntitysCreatures) MaxDamageChanges(creature models.Creatures) (name string, cant int) {

	if creature.Ice >= cant {
		name = "Ice"
		cant = creature.Ice
	}
	if creature.Physical >= cant {
		name = "Physical"
		cant = creature.Physical
	}
	if creature.Earth >= cant {
		name = "Earth"
		cant = creature.Earth
	}
	if creature.Fire >= cant {
		name = "Fire"
		cant = creature.Fire
	}
	if creature.Death >= cant {
		name = "Death"
		cant = creature.Death
	}
	if creature.Energy >= cant {
		name = "Energy"
		cant = creature.Energy
	}
	if creature.Holy >= cant {
		name = "Holy"
		cant = creature.Holy
	}

	return
}

func (ec *EntitysCreatures) MaxDamageChangesBoss(creature models.Bosses) (name string, cant int) {

	if creature.Ice >= cant {
		name = "Ice"
		cant = 100 - creature.Ice
	}
	if creature.Holy >= cant {
		name = "Holy"
		cant = 100 - creature.Holy
	}
	if creature.Earth >= cant {
		name = "Earth"
		cant = 100 - creature.Earth
	}
	if creature.Fire >= cant {
		name = "Fire"
		cant = 100 - creature.Fire
	}
	if creature.Death >= cant {
		name = "Death"
		cant = 100 - creature.Death
	}
	if creature.Energy >= cant {
		name = "Energy"
		cant = 100 - creature.Energy
	}
	if creature.Physical >= cant {
		name = "Physical"
		cant = 100 - creature.Physical
	}

	return
}
