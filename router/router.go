package router

import (
	"net/http"
	"text/template"

	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/handler"
	"github.com/Mau005/KraynoSerer/middleware"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()
	router.Use(middleware.CommonMiddleware)
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fileImage := http.FileServer(http.Dir("./data/image"))
	router.PathPrefix("/image/").Handler(http.StripPrefix("/image/", fileImage))

	var HomeHandler handler.HomeHandler
	router.HandleFunc("/", HomeHandler.Home).Methods("GET")
	router.HandleFunc("/todays", HomeHandler.Todays).Methods("GET")
	router.HandleFunc("/todays_post/{id}", HomeHandler.TodaysPost).Methods("GET")

	router.HandleFunc("/404", handler.Page404).Methods("GET")
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/test.html")
		if err != nil {
			return
		}
		var api controller.ApiController
		sc := api.GetBaseWeb(r)
		template.Execute(w, sc)
	})

	var AccountHandler handler.AccountHandler
	router.HandleFunc("/create_user", AccountHandler.CreateAccount).Methods("POST")
	router.HandleFunc("/login", AccountHandler.Login).Methods("POST")
	router.HandleFunc("/News", nil).Methods("GET")
	router.HandleFunc("/Todays", nil).Methods("GET")
	router.HandleFunc("/Todays/{page}", nil).Methods("GET")
	router.HandleFunc("/logout", AccountHandler.Logout).Methods("GET")

	security := router.PathPrefix("/auth").Subrouter()
	security.Use(middleware.CommonMiddleware)
	security.Use(middleware.SessionMiddleware)

	var imageHanlder handler.ImageHandler
	security.HandleFunc("/upload_image", imageHanlder.UploadImageHandler).Methods("GET")
	security.HandleFunc("/upload_image", imageHanlder.LoadImage).Methods("POST")

	security.HandleFunc("/my_profile", AccountHandler.MyProfileHandler).Methods("GET")
	security.HandleFunc("/my_setting", AccountHandler.MyProfileSettingPOST).Methods("POST")
	security.HandleFunc("/change_password", AccountHandler.MyProfileChangePasswordHandler).Methods("POST")
	security.HandleFunc("/add_comment", AccountHandler.AddCommentTodays).Methods("POST")

	var adminHandler handler.AdminHandler
	security.HandleFunc("/admin", adminHandler.Lobby).Methods("GET")
	security.HandleFunc("/todays_aproved", adminHandler.TodaysAproved).Methods("GET")
	security.HandleFunc("/todays_aproved/{id}", adminHandler.TodaysAprovedPOST).Methods("POST")

	return router
}
