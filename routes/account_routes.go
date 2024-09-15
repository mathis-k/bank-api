package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
)

func RegisterAccountRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api/accounts").Subrouter()
	subRouter.HandleFunc("", controllers.GetAccounts).Methods("GET")
	subRouter.HandleFunc("{id}", controllers.GetAccountById).Methods("GET")
	subRouter.HandleFunc("", controllers.CreateAccount).Methods("POST")
	subRouter.HandleFunc("{id}", controllers.DeleteAccount).Methods("DELETE")
}
