package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/middleware"
	"github.com/mathis-k/bank-api/models"
	"github.com/mathis-k/bank-api/routes"
	"log"
	"net/http"
	"os"
)

type APIServer struct {
	listenAddress string
	database      *models.DB
}
type APIResponse struct {
	Message string `json:"message"`
}

func jsonMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(APIResponse{message}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewAPIServer() *APIServer {

	if err := godotenv.Load(); err != nil {
		log.Println("✖ No .env file found")
		return &APIServer{}
	}
	listenAddress := os.Getenv("API_SERVER_ADDRESS")
	if listenAddress == "" {
		log.Fatal("✖ API_SERVER_ADDRESS environment variable not set")
	}
	database := &models.DB{}

	log.Println("✔ New API server created on address:", listenAddress)
	if err := database.Connect(); err != nil {
		panic("✖ Could not connect to database")
	}
	return &APIServer{
		listenAddress: listenAddress,
		database:      database,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api", s.handleStartPage).Methods(http.MethodGet)
	router.Use(middleware.AuthMiddleware)

	routes.RegisterUserRoutes(router)
	routes.RegisterAccountRoutes(router)
	routes.RegisterTransactionRoutes(router)
	routes.RegisterAuthRoutes(router)

	log.Printf("✔ API server is running on localhost%s/ ... 🚀", s.listenAddress)
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		log.Println("⚠ Error whilst listening:", err)
		log.Println("⚠ ... Shutting down server ...")
		s.Shutdown()
		return
	}
}
func (s *APIServer) Shutdown() {
	if err := s.database.Disconnect(); err != nil {
		log.Printf("⚠ Error disconnecting from database: %v", err)
	}
	log.Println("✔ API server has been shut down.")
}

func (s *APIServer) handleStartPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("https://github.com/mathis-k/bank-api")); err != nil {
		return
	}
}
