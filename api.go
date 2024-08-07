package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type APIServer struct {
	listenAddress string
}

func NewAPI(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		panic(fmt.Errorf("%w: Error starting server: %v", err, s.listenAddress))
	}
	log.Println("API server is running on", s.listenAddress)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	case http.MethodGet:
		return s.handleGetAccounts(w, r)
	case http.MethodDelete:
		return s.handleDeleteAccount(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}
}
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
