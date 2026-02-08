package models

type StructModel struct {
	Lenguaje        map[string]string
	Daily           []string
	Email           string
	UserName        string
	Access          uint8
	TitleError      string
	MSGError        string
	RouterError     string
	NameButtonError string
	StreamMode      bool
	LenguajeDefault string
}

type FrontEnd struct {
	Title     string
	Meta      string
	Link      string
	Preload   string
	Footer    string
	ButtonVol string
	Login     string
	Scripts   string
}
