package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/db"
	"github.com/mathis-k/bank-api/models"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddress string
	database      models.Database
}
type APIResponse struct {
	Message string `json:"message"`
}

func json_message(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(APIResponse{message}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewAPIServer(listenAddress string) *APIServer {
	database := &db.MongoDB{}
	if err := database.Connect(); err != nil {
		panic("Could not connect to database")
	}

	return &APIServer{
		listenAddress: listenAddress,
		database:      database,
	}
}
func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/", s.handleStartPage).
		Methods(http.MethodGet)
	router.HandleFunc("/account", s.handleAccount).
		Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/account/{id}", s.handleAccountByID).
		Methods(http.MethodGet, http.MethodDelete, http.MethodPut)

	log.Printf("API server is running on localhost%s ... 🚀", s.listenAddress)
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		log.Panicf("✖ %w: Error starting server: %v", err, s.listenAddress)
	}
}
func (s *APIServer) handleStartPage(w http.ResponseWriter, r *http.Request) {
	// Set the content type to plain text
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(welcomeMessage)); err != nil {
		return
	}
}
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleCreateAccount(w, r)
	case http.MethodGet:
		s.handleGetAccounts(w, r)
	default:
		json_message(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}
func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetAccountByID(w, r)
	case http.MethodDelete:
		s.handleDeleteAccount(w, r)
	case http.MethodPut:
		json_message(w, http.StatusNotImplemented, "Method not implemented yet")
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {

	accounts, err := s.database.GetAllAccounts()
	if err != nil {
		json_message(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(accounts) == 0 {
		json_message(w, http.StatusNoContent, "No accounts found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		json_message(w, http.StatusInternalServerError, err.Error())
	}
}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		json_message(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		json_message(w, http.StatusInternalServerError, err.Error())
	}
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		json_message(w, http.StatusBadRequest, err.Error())
		return
	}
	json_message(w, http.StatusOK, fmt.Sprintf("Account with ID: %d deleted successfully!", id))
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json_message(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.FirstName == "" {
		json_message(w, http.StatusBadRequest, "Missing Firstname")
		return
	} else if req.LastName == "" {
		json_message(w, http.StatusBadRequest, "Missing Lastname")
		return
	} else if req.Email == "" {
		json_message(w, http.StatusBadRequest, "Missing Email")
		return
	}

	account, err := models.NewAccount(&req)
	if err != nil {
		json_message(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.database.CreateAccount(account); err != nil {
		if err.Error() == fmt.Sprintf("an account with the email %s already exists", req.Email) {
			json_message(w, http.StatusConflict, err.Error())
		} else {
			json_message(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		json_message(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	return
}

const welcomeMessage = `Welcome to the Bank JSON API Server! :)
🏦
Available endpoints:
GET /account - get all accounts
POST /account - create a new account
GET /account/{id} - get account by ID
PUT /account/{id} - update account by ID
DELETE /account/{id} - delete account by ID`
