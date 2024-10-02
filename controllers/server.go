package controllers

import (
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/models"
	"log"
	"net/http"
	"os"
)

type APIServer struct {
	ListenAddress string
	Database      *models.DB
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
		ListenAddress: listenAddress,
		Database:      database,
	}
}

func (s *APIServer) HandleStartPage(w http.ResponseWriter, r *http.Request) {
	r.Context()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("https://github.com/mathis-k/bank-api")); err != nil {
		return
	}
}
