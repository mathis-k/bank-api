package routes

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
	"github.com/mathis-k/bank-api/middleware"
)

func RegisterUserRoutes(router *mux.Router, controllers *controllers.APIServer) {
	subRouter := router.PathPrefix("/api/user").Subrouter()
	subRouter.Use(middleware.AuthMiddleware)
	subRouter.HandleFunc("", controllers.GetUser).Methods("GET")
	subRouter.HandleFunc("", controllers.UpdateUser).Methods("PUT")
}
