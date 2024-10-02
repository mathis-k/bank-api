package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
)

func RegisterAuthRoutes(router *mux.Router, controllers *controllers.APIServer) {
	subRouter := router.PathPrefix("/api/auth").Subrouter()
	subRouter.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	subRouter.HandleFunc("/login", controllers.LoginUser).Methods("POST")
}
