package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
	"github.com/mathis-k/bank-api/middleware"
)

func RegisterTransactionRoutes(router *mux.Router, controllers *controllers.APIServer) {
	subRouter := router.PathPrefix("/api/transactions").Subrouter()
	subRouter.Use(middleware.AuthMiddleware)
	subRouter.HandleFunc("", controllers.GetTransactions).Methods("GET")
	subRouter.HandleFunc("/{id}", controllers.GetTransactionById).Methods("GET")

	subsubRouter := subRouter.PathPrefix("/account").Subrouter()
	subsubRouter.Use(controllers.Database.CheckAccountPermissionMiddleware)
	subsubRouter.HandleFunc("/{number}", controllers.GetTransactionsFromAccount).Methods("GET")
	subsubRouter.HandleFunc("/{number}/deposit", controllers.DepositToAccount).Methods("POST")
	subsubRouter.HandleFunc("/{number}/withdraw", controllers.WithdrawFromAccount).Methods("POST")
	subsubRouter.HandleFunc("/{number}/transfer", controllers.TransferBetweenAccounts).Methods("POST")
}
