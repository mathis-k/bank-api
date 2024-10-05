package main

import (
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/controllers"
	"github.com/mathis-k/bank-api/routes"
	"log"
	"net/http"
)

func Shutdown(s *controllers.APIServer) {
	if err := s.Database.Disconnect(); err != nil {
		log.Printf("⚠ Error disconnecting from database: %v", err)
	}
	log.Println("✔ API server has been shut down.")
}
func Run(s *controllers.APIServer) {
	router := mux.NewRouter()
	router.HandleFunc("/api", s.HandleStartPage).Methods(http.MethodGet)

	routes.RegisterUserRoutes(router, s)
	routes.RegisterAccountRoutes(router, s)
	routes.RegisterTransactionRoutes(router, s)
	routes.RegisterAuthRoutes(router, s)

	log.Printf("✔ API server is running on localhost%s/ ... 🚀", s.ListenAddress)
	err := http.ListenAndServe(s.ListenAddress, router)
	if err != nil {
		log.Println("⚠ Error whilst listening:", err)
		log.Println("⚠ ... Shutting down server ...")
		Shutdown(s)
		return
	}
}
func main() {
	api := controllers.NewAPIServer()
	Run(api)
}
