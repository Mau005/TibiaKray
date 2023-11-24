package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Email      string
	UserName   string
	Access     uint8
	StreamMode bool
	Lenguaje   string
	jwt.StandardClaims
}
