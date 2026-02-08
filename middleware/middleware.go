package middleware

import (
	"net/http"

	"github.com/Mau005/KraynoSerer/configuration"
	"github.com/Mau005/KraynoSerer/controller"
	"github.com/Mau005/KraynoSerer/handler"
	"github.com/gorilla/context"
)

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, thorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := configuration.Store.Get(r, configuration.NAME_SESSION)
		if err != nil {
			http.Error(w, "Error al obtener la sesión", http.StatusInternalServerError)
			return
		}
		var api controller.ApiController
		if tokenStr, ok := session.Values["token"].(string); !ok {
			var api controller.ApiController
			sc := api.GetBaseWeb(r)
			var ErrorHandler handler.ErrorHandler
			ErrorHandler.PageErrorMSG(http.StatusNetworkAuthenticationRequired, configuration.NotAuthorized, configuration.ROUTER_INDEX, w, r, sc)
			return
		} else {
			err = api.AuthenticateJWT(tokenStr)
			if err != nil {
				var api controller.ApiController
				sc := api.GetBaseWeb(r)
				api.SaveSession(nil, w, r) //cerramos la secion
				var ErrorHandler handler.ErrorHandler
				ErrorHandler.PageErrorMSG(http.StatusUnauthorized, configuration.ErrorExpireSession, configuration.ROUTER_INDEX, w, r, sc)
				return
			}

		}
		context.Set(r)
		next.ServeHTTP(w, r)

	})
}
