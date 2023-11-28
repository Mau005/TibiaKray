package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/gorilla/mux"
)

type VotedHandler struct{}

func (vt *VotedHandler) AddVotedTodays(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)

	idTodays, err := strconv.ParseUint(args["id"], 10, 64)
	if err != nil {
		return
	}
	var api controller.ApiController
	sc := api.GetBaseWeb(r)
	var votedManager controller.VotedController
	_, err = votedManager.TodaysVoted(uint(idTodays), sc)
	if err != nil {
		var errorHanlder ErrorHandler
		errorHanlder.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorPrivileges, fmt.Sprintf(configuration.ROUTER_TODAYS_POST, idTodays), w, r, sc)
		return
	}
}
