package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mathis-k/bank-api/db"
	"github.com/mathis-k/bank-api/models"
	"log"
	"math"
	"net/http"
	"strconv"
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error
type APIServer struct {
	listenAddress string
	database      models.Database
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
	router.HandleFunc("/transfer", s.handleTransfer).
		Methods(http.MethodPost)

	log.Printf("✔ API server is running on localhost%s/ ... 🚀", s.listenAddress)
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		log.Printf("✖ %w: Error starting server: %v", err, s.listenAddress)
		log.Println("... Shutting down server ...")
		return
	}
}
func (s *APIServer) handleStartPage(w http.ResponseWriter, r *http.Request) {

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
		jsonMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
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
		s.handleUpdateAccount(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	maxResults := uint64(math.MaxUint64)
	if r.URL.Query().Get("maxResult") != "" {
		input := r.URL.Query().Get("maxResult")
		val, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid max query parameter '%s': %v", input, err)
			jsonMessage(w, http.StatusBadRequest, msg)
			return
		}
		maxResults = val
	}

	accounts, err := s.database.GetAllAccounts(maxResults)
	if err != nil {
		jsonMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(accounts) == 0 {
		jsonMessage(w, http.StatusNoContent, "No accounts found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		jsonMessage(w, http.StatusInternalServerError, err.Error())
	}
}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	account, err := s.database.GetAccountByID(id)
	if err != nil {
		if err.Error() == db.NoAccountFound {
			jsonMessage(w, http.StatusNotFound, err.Error())
		} else if err.Error() == db.InvalidID {
			jsonMessage(w, http.StatusBadRequest, err.Error())
		} else {
			jsonMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		jsonMessage(w, http.StatusInternalServerError, err.Error())
	}
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req models.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := s.database.UpdateAccount(id, &req)
	if err != nil {
		if err.Error() == db.NoAccountFound {
			jsonMessage(w, http.StatusNotFound, err.Error())
		} else if err.Error() == db.InvalidID || err.Error() == models.InvalidEmail || err.Error() == models.InvalidFirstName || err.Error() == models.InvalidLastName {
			jsonMessage(w, http.StatusBadRequest, err.Error())
		} else if err.Error() == fmt.Sprintf("an account with the email %s already exists, please choose another email", req.Email) {
			jsonMessage(w, http.StatusConflict, err.Error())
		} else {
			jsonMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		jsonMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	account, err := s.database.DeleteAccount(id)
	if err != nil {
		if err.Error() == db.NoAccountFound {
			jsonMessage(w, http.StatusNotFound, err.Error())
		} else if err.Error() == db.InvalidID {
			jsonMessage(w, http.StatusBadRequest, err.Error())
		} else {
			jsonMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		jsonMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.FirstName == "" {
		jsonMessage(w, http.StatusBadRequest, "Missing Firstname")
		return
	} else if req.LastName == "" {
		jsonMessage(w, http.StatusBadRequest, "Missing Lastname")
		return
	} else if req.Email == "" {
		jsonMessage(w, http.StatusBadRequest, "Missing Email")
		return
	}
	account, err := s.database.CreateAccount(&req)
	if err != nil {
		if err.Error() == fmt.Sprintf("an account with the email %s already exists", req.Email) {
			jsonMessage(w, http.StatusConflict, err.Error())
		} else {
			jsonMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		jsonMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {

	return
}

const welcomeMessage = `Welcome to the Bank JSON API Server! :)

Available endpoints:
GET /account - get all accounts
POST /account - create a new account
GET /account/{id} - get account by ID
PUT /account/{id} - update account by ID
DELETE /account/{id} - delete account by ID
PUT /transfer - transfer money between accounts`
