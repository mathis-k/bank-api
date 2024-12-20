package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
	"github.com/mathis-k/bank-api/middleware"
)

func RegisterAccountRoutes(router *mux.Router, controllers *controllers.APIServer) {
	subRouter := router.PathPrefix("/api/accounts").Subrouter()
	subRouter.Use(middleware.AuthMiddleware)
	subRouter.HandleFunc("", controllers.GetAccounts).Methods("GET")
	subRouter.HandleFunc("", controllers.CreateAccount).Methods("POST")

	subsubRouter := subRouter.PathPrefix("/{number}").Subrouter()
	subsubRouter.Use(controllers.Database.CheckAccountPermissionMiddleware)
	subsubRouter.HandleFunc("", controllers.GetAccountByNumber).Methods("GET")
	subsubRouter.HandleFunc("", controllers.DeleteAccount).Methods("DELETE")
}
