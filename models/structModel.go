package models

type StructModel struct {
	Lenguaje        map[string]string
	News            []News
	NewsTicket      []NewsTicket
	Daily           []string
	Email           string
	UserName        string
	Group           uint8
	TitleError      string
	MSGError        string
	RouterError     string
	NameButtonError string
	StreamMode      bool
}
