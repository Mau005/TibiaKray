package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
	var votedManager controller.VotedController
	voted := votedManager.TodaysVoted(uint(idTodays))
	fmt.Println(voted.Status)
}
