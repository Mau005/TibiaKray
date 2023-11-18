package controller

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var Lenguaje map[string]map[string]string = make(map[string]map[string]string) //load data

type ApiController struct{}

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
		content := strings.Split(record[0], ";")
		if content[0] == "base" {
			lenguajeList = content[0:]
			for _, value := range content[0:] {
				Lenguaje[value] = make(map[string]string)
			}
		} else {
			for index, value := range lenguajeList {
				//fmt.Println(fmt.Sprintf("Idioma: %s Variable: %s Valor: %s", value, content[0], content[index]))
				Lenguaje[value][content[0]] = content[index]
			}
		}
	}
	log.Println(fmt.Sprintf("[Lenguaje] Multi language has been started supporting: %v", lenguajeList))
	return nil
}

func (ac *ApiController) CompareCryptPassword(password, passwordTwo string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordTwo))
}

func (ac *ApiController) generateClaims(account *models.Account) *models.Claims {
	expirationTime := time.Now().Add(configuration.EXPIRATION_TOKEN * time.Hour)
	return &models.Claims{
		Email:    account.Email,
		UserName: account.Name,
		Group:    account.Groups,
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

func (ac *ApiController) GetBaseWeb(r *http.Request) (sc models.StructModel) {
	var api ApiController
	claims, _ := api.GetSessionClaims(r)

	//Configuration User Session
	sc.UserName = claims.UserName
	sc.Email = claims.Email
	sc.Group = claims.Group
	//Configuration User Session
	sc.Lenguaje = Lenguaje["en"]
	var accController AccountController
	account, err := accController.GetAccount(sc.Email)
	if err == nil {
		sc.Lenguaje = Lenguaje[account.Languaje]
	}

	//Config Lenguaje

	//End Lenguaje

	sc.News = *News
	sc.NewsTicket = *NewsTicket
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

func (ac *ApiController) GenerateEncrypt()
