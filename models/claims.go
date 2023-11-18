package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Email    string
	UserName string
	Group    uint8
	jwt.StandardClaims
}
