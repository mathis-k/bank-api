package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

/*
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
*/

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}
func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/account", s.handleAccount).
		Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/account/{id}", s.handleAccountByID).
		Methods(http.MethodGet, http.MethodDelete, http.MethodPut)

	log.Printf("API server is running on localhost%s ... 🚀", s.listenAddress)
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		panic(fmt.Errorf("%w: Error starting server: %v", err, s.listenAddress))
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleCreateAccount(w, r)
	case http.MethodGet:
		s.handleGetAccounts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		http.Error(w, "Method not implemented yet", http.StatusNotImplemented)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := []Account{
		*NewAccount("John", "Doe", "john.doe@example.com"),
		*NewAccount("Jane", "Doe", "jane.doe@example.com"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Printf("ID: %d\n", id)
	account := NewAccount("John", "Doe", "john.doe@example.com")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Account with ID: %d deleted successfully!", id)
	if err := json.NewEncoder(w).Encode(APIResponse{msg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	return
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	return
}
