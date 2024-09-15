package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
)

func RegisterTransactionRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api/transactions").Subrouter()
	subRouter.HandleFunc("", controllers.GetTransactions).Methods("GET")
	subRouter.HandleFunc("/{id}", controllers.GetTransactionById).Methods("GET")
	subRouter.HandleFunc("/account/{id}", controllers.GetTransactionsFromAccount).Methods("GET")
	subRouter.HandleFunc("/account/{id}/deposit", controllers.DepositToAccount).Methods("POST")
	subRouter.HandleFunc("/account/{id}/withdraw", controllers.WithdrawFromAccount).Methods("POST")
	subRouter.HandleFunc("/account/{id}/transfer", controllers.TransferBetweenAccounts).Methods("POST")
}
