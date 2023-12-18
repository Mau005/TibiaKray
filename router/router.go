package router

import (
	"net/http"

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

	var imageHandler handler.ImageHandler
	router.HandleFunc("/todays/{page}", imageHandler.GetPhotosHandler).Methods("GET")

	router.HandleFunc("/404", handler.Page404).Methods("GET")
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var monsterContr controller.EntitysCreatures
		monsterContr.CollectorCreature()
	})

	var AccountHandler handler.AccountHandler
	router.HandleFunc("/create_user", AccountHandler.CreateAccount).Methods("POST")
	router.HandleFunc("/login", AccountHandler.Login).Methods("POST")
	router.HandleFunc("/News", nil).Methods("GET")
	var errHandler handler.ErrorHandler
	router.NotFoundHandler = http.HandlerFunc(errHandler.DefaultError)
	router.HandleFunc("/logout", AccountHandler.Logout).Methods("GET")

	var recoveryHandler handler.RecoveryHandler
	router.HandleFunc("/recovery/{code}", recoveryHandler.RecoveryHandler).Methods("GET")
	router.HandleFunc("/recovery", recoveryHandler.RecoveryAccount).Methods("POST")
	router.HandleFunc("/recovery_account", recoveryHandler.RecoveryChangePassword).Methods("POST")

	var toolshandler handler.ToolsHandler
	router.HandleFunc("/shared_loot", toolshandler.SharedLootHandler).Methods("GET")
	router.HandleFunc("/shared_loot", toolshandler.SharedLootProcess).Methods("POST")
	router.HandleFunc("/sharedexp", toolshandler.SharedExpHanlder).Methods("GET")
	router.HandleFunc("/sharedexp", toolshandler.SharedExpProcess).Methods("POST")
	router.HandleFunc("/tools", toolshandler.ToolsHandlerItems).Methods("GET")

	var creaturesHandler handler.CreaturesHandler
	router.HandleFunc("/creatures", creaturesHandler.CreaturesHandler).Methods("GET")
	router.HandleFunc("/bosses", creaturesHandler.BossesHandler).Methods("GET")

	security := router.PathPrefix("/auth").Subrouter()
	security.Use(middleware.CommonMiddleware)
	security.Use(middleware.SessionMiddleware)

	var imageHanlder handler.ImageHandler
	security.HandleFunc("/upload_image", imageHanlder.UploadHandlerDefault).Methods("GET")
	security.HandleFunc("/decisions_upload", imageHandler.UploadDessioningHandler).Methods("POST")
	security.HandleFunc("/upload_image", imageHanlder.LoadImage).Methods("POST")
	security.HandleFunc("/upload_url", imageHandler.UploadUrl).Methods("POST")

	security.HandleFunc("/my_profile", AccountHandler.MyProfileHandler).Methods("GET")
	security.HandleFunc("/my_setting", AccountHandler.MyProfileSettingPOST).Methods("POST")
	security.HandleFunc("/change_password", AccountHandler.MyProfileChangePasswordHandler).Methods("POST")
	security.HandleFunc("/add_comment", AccountHandler.AddCommentTodays).Methods("POST")

	security.HandleFunc("/my_favorite_picture", AccountHandler.MyProfilePictureHandler).Methods("GET")
	security.HandleFunc("/my_players", AccountHandler.MyProfilePLayers).Methods("GET")
	security.HandleFunc("/search_player", AccountHandler.SearchMyPlayer).Methods("POST")

	var votedHanlder handler.VotedHandler
	security.HandleFunc("/voted_todays/{id}", votedHanlder.AddVotedTodays).Methods("POST")

	var adminHandler handler.AdminHandler
	security.HandleFunc("/admin", adminHandler.Lobby).Methods("GET")
	security.HandleFunc("/todays_aproved", adminHandler.TodaysAproved).Methods("GET")
	security.HandleFunc("/todays_aproved/{id}", adminHandler.TodaysAprovedPOST).Methods("POST")
	security.HandleFunc("/user_register", adminHandler.UserRegisterHandler).Methods("GET")

	return router
}
