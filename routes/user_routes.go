package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
)

func RegisterUserRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api/user").Subrouter()
	subRouter.HandleFunc("", controllers.GetUser).Methods("GET")
	subRouter.HandleFunc("", controllers.UpdateUser).Methods("PUT")
	subRouter.HandleFunc("/accounts", controllers.GetAccountsFromUser).Methods("GET")
	subRouter.HandleFunc("/transactions", controllers.GetTransactionsFromUser).Methods("GET")
}
