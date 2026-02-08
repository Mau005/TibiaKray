package controller

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Mau005/KraynoSerer/components"
	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/database"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var Lenguaje map[string]map[string]string = make(map[string]map[string]string)   //load data
var LenguajeInternal map[string]map[int]string = make(map[string]map[int]string) //load data
var Manager *ManagerController
var CountVisit uint
var RecoveryAccount map[string]RecoveryController //sistema para gestionar las recuperaciones

var ExtencionImage = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"gif":  true,
	"bmp":  true,
}

type ApiController struct {
	Components components.Components
}

func (ac *ApiController) InitServices() error {
	RecoveryAccount = make(map[string]RecoveryController) //Iniciamos la recuperacion de cuentas
	CountVisit = 0
	//Load COnfig
	err := configuration.LoadConfiguration(configuration.PATH_CONFIG)
	if err != nil {
		log.Println(err)
		return err
	}
	//End Load Config
	//Load Database
	err = database.ConnectionDataBase()
	if err != nil {
		log.Println(err)
		return err
	}
	//End Load DataBase
	//Init Lenguaje
	err = ac.InitLenguaje(configuration.PATH_LENGUAJE_CLIENT)
	if err != nil {
		log.Println(err)
		return err
	}
	err = ac.InitLenguajeServer(configuration.PATH_LENGUAJE_SERVER)
	if err != nil {
		log.Println(err)
		return err
	}
	//End Lenguaje

	//init News for tibia.com
	//End News for tibia.com

	//init rashid
	var toolsManager ToolsController
	toolsManager.InitRashid()
	//end rashid

	//init default reset web for server tibia
	go ac.ResetDefaultWeb() //Iniciamos los servicios cada 24 hrs
	//end init
	Manager, err = NewManagerController()
	if err != nil {
		log.Println("Error al cargar Manager Controller" + err.Error())
		return err
	}
	go Manager.Update()

	return nil

}

func (ac *ApiController) GenerateCryptPassword(password string) string {
	hasedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hasedPassword)
}

func (ac *ApiController) InitLenguaje(path_file string) error {

	file, err := os.Open(path_file)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lenguajeList := []string{}
	for {
		record, err := reader.Read()

		if err != nil {
			break
		}
		if record[0] == "base" {
			lenguajeList = record[0:]
			for _, value := range record[0:] {
				Lenguaje[value] = make(map[string]string)
			}
		} else {
			for index, value := range lenguajeList {
				//fmt.Println(fmt.Sprintf("Idioma: %s Variable: %s Valor: %s", value, content[0], content[index]))
				Lenguaje[value][record[0]] = record[index]
			}
		}
	}
	log.Println(fmt.Sprintf("[Lenguaje] Multi language has been started supporting: %v", lenguajeList))
	return nil
}

func (ac *ApiController) InitLenguajeServer(path_file string) error {

	file, err := os.Open(path_file)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	list_map := make([]string, 0, 5)
	for {
		record, e := reader.Read()
		if e != nil {
			break
		}

		if record[0] == "base" {
			for _, value := range record {
				LenguajeInternal[value] = make(map[int]string)
				list_map = append(list_map, value)
			}
			continue
		}

		for index, value := range list_map {
			//index, valor es, en
			idVar, err := strconv.ParseInt(record[0], 10, 8)
			if err != nil {
				log.Println(err, idVar, record[0])
				continue
			}
			//fmt.Println(fmt.Sprintf("KEY: |%d| Value: |%s|", int(idVar), record[index]))
			LenguajeInternal[value][int(idVar)] = record[index]
		}
	}

	return nil
}

func (ac *ApiController) CompareCryptPassword(password, passwordTwo string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordTwo))
}

func (ac *ApiController) generateClaims(account *models.Account) *models.Claims {
	expirationTime := time.Now().Add(configuration.EXPIRATION_TOKEN * time.Hour)
	return &models.Claims{
		Email:      account.Email,
		UserName:   account.Name,
		Access:     account.Access,
		StreamMode: account.StreamMode,
		Lenguaje:   account.Languaje,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
}

func (ac *ApiController) GenerateToken(account *models.Account) (tokenString string, err error) {

	claims := ac.generateClaims(account)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(configuration.Security))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func (ac *ApiController) SaveSession(tokenString *string, w http.ResponseWriter, r *http.Request) {
	session, _ := configuration.Store.Get(r, configuration.NAME_SESSION)
	if tokenString == nil {
		session.Values["token"] = nil
	} else {
		session.Values["token"] = *tokenString
	}

	session.Save(r, w)
}

func (ac *ApiController) AuthenticateJWT(tokenSession string) error {

	token, err := jwt.Parse(tokenSession, func(token *jwt.Token) (interface{}, error) {
		return []byte(configuration.Security), nil
	})

	if err != nil || !token.Valid {
		return err
	}

	return nil

}

func (ac *ApiController) GetSessionClaims(r *http.Request) (*models.Claims, error) {
	claims := &models.Claims{}
	session, err := configuration.Store.Get(r, configuration.NAME_SESSION)
	if err != nil {
		return claims, err
	}

	token, ok := session.Values["token"].(string)
	if !ok {
		return claims, errors.New("Token de session invalido")
	}
	tokenKey := []byte(configuration.Security)
	tokenParser := jwt.Parser{}

	_, err = tokenParser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenKey, nil
	})
	if err != nil {
		return claims, err
	}
	return claims, nil
}

func (ac *ApiController) GetFronEnd(title, description string, claim models.Claims, lenguaje map[string]string) (front models.FrontEnd) {

	front.Title = ac.Components.CrceateTitle(title)
	front.Meta = ac.Components.CreateMetaDefault(title, description)
	front.Link = ac.Components.CreateLink()
	front.ButtonVol = ac.Components.CreateButtonVolumen()
	front.Login = ac.Components.CreateLogin(lenguaje)
	return
}

func (ac *ApiController) GetBaseWeb(r *http.Request) (sc models.StructModel) {
	var api ApiController
	claims, _ := api.GetSessionClaims(r)

	//Configuration User Session
	sc.UserName = claims.UserName
	sc.Email = claims.Email
	sc.Access = claims.Access
	sc.StreamMode = claims.StreamMode

	if claims.Lenguaje == "" {
		sc.LenguajeDefault = configuration.DEFAULT_LENGUAJE
	} else {
		sc.LenguajeDefault = claims.Lenguaje
	}

	//Configuration User Session
	sc.Lenguaje = Lenguaje[configuration.DEFAULT_LENGUAJE] //posibilidad de usar las ips para identificar el pais para el idioma!
	var accController AccountController
	account, err := accController.GetAccount(sc.Email)
	if err == nil {
		sc.Lenguaje = Lenguaje[account.Languaje]
	}

	//Config Lenguaje

	return sc
}

func (ac *ApiController) GetWebError(r *http.Request, title, content, router, namebutton string) (sc models.StructModel) {
	sc = ac.GetBaseWeb(r)
	sc.MSGError = content
	sc.TitleError = title
	sc.RouterError = router
	sc.NameButtonError = namebutton
	return sc
}

func (ac *ApiController) GetSessionAccount(r *http.Request) (acc *models.Account, err error) {
	claims, err := ac.GetSessionClaims(r)
	if err != nil {
		return acc, err
	}
	var accContr AccountController

	acc, err = accContr.GetAccount(claims.Email)
	if err != nil {
		return acc, err
	}

	return acc, nil
}

func (ac *ApiController) GenerateHash(content string) string {
	h := sha256.New()
	io.WriteString(h, content)
	return hex.EncodeToString(h.Sum(nil))
}

func (ac *ApiController) GenerateEncrypt(fileOrigin multipart.File, origin, nameFile string, account *models.Account, idEncrypt uint, validations map[string]bool) {
	defer fileOrigin.Close()
	verify := strings.Split(nameFile, ".")

	if !(validations[verify[len(verify)-1]]) {
		log.Println("La extencion indicada no es correcta")
		return
	}

	EncrName := ac.GenerateHash(fmt.Sprintf("%d", idEncrypt))
	EncrNameExtencion := EncrName + "." + verify[len(verify)-1]

	pathOrigin := origin + EncrName + EncrNameExtencion
	pathDir := origin
	fmt.Println(pathDir, pathOrigin)

}

func (ac *ApiController) NormalizeString(count int, content string) (result string) {

	if !(len(content) >= count) {
		return content
	}

	for _, characteres := range content {
		result += string(characteres)
		if count <= len(result) {
			result += "..."
			break
		}
	}

	return result
}

func (ac *ApiController) GenerateUUid() string {
	//Se generan los UUID para los personajes registrados
	id := uuid.New()

	return id.String()
}

func (ac *ApiController) ResetDefaultWeb() {
	for {
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

		var toolsManager ToolsController
		toolsManager.InitRashid()

		//reset sharedloot
		configuration.SharedLootHightNow = models.SharedLoot{}

	}

}
func (ac *ApiController) DownloadImage(url, filePath string) error {
	pathOrigin := configuration.PATH_STATIC_PUBLIC + filePath

	// Si el archivo ya existe, salir
	if _, err := os.Stat(pathOrigin); err == nil {
		return errors.New("el archivo ya existe")
	}

	// 🔹 Crear directorio si no existe
	dir := filepath.Dir(pathOrigin)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Println("error creando directorio:", err)
		return err
	}

	// Descargar imagen
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	// Validar HTTP status
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error descargando imagen: status %d", response.StatusCode)
	}

	// Crear archivo
	file, err := os.Create(pathOrigin)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	// Copiar contenido
	if _, err = io.Copy(file, response.Body); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Imagen descargada con éxito: %s\n", filePath)
	return nil
}
