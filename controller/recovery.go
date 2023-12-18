package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mau005/KraynoSerer/configuration"
)

type RecoveryController struct {
	AccName  string
	Email    string
	Hash     string
	Code     string
	ExpireAt time.Time
}

func (rc *RecoveryController) NewRecoveryAccount(accname, email, lenguaje string) (RecoveryController, error) {
	if accname == "" && email == "" {
		return RecoveryController{}, errors.New(LenguajeInternal[lenguaje][configuration.ErrorEmptyField])
	}

	var api ApiController
	hash := api.GenerateHash(api.GenerateUUid())

	now := time.Now().Add(5 * time.Minute)
	return RecoveryController{AccName: accname, Email: email, ExpireAt: now, Hash: hash, Code: fmt.Sprintf("%d", rand.Intn(999999))}, nil

}
