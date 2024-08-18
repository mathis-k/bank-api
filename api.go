package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/auth"
	"github.com/mathis-k/bank-api/db"
	"github.com/mathis-k/bank-api/models"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

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

func NewAPIServer() *APIServer {

	if err := godotenv.Load(); err != nil {
		log.Println("✖ No .env file found")
		return &APIServer{}
	}
	listenAddress := os.Getenv("API_SERVER_ADDRESS")
	if listenAddress == "" {
		log.Fatal("✖ API_SERVER_ADDRESS environment variable not set")
	}
	database := &db.MongoDB{}

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
	router.HandleFunc("/", s.handleStartPage).
		Methods(http.MethodGet)
	router.HandleFunc("/account", withJWTAuth(s.handleAccount)).
		Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/account/{id}", withJWTAuth(s.handleAccountByID)).
		Methods(http.MethodGet, http.MethodDelete, http.MethodPut)
	router.HandleFunc("/transfer", withJWTAuth(s.handleTransfer)).
		Methods(http.MethodPost)

	log.Printf("✔ API server is running on localhost%s/ ... 🚀", s.listenAddress)
	auth.GenerateAdminJWT()
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

func withJWTAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonMessage(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := auth.VerifyJWT(tokenString)
		if err != nil {
			if err.Error() == auth.TOKEN_EXPIRED || err.Error() == auth.INVALID_TOKEN {
				jsonMessage(w, http.StatusUnauthorized, err.Error())
			} else {
				jsonMessage(w, http.StatusBadRequest, err.Error())
			}
			return
		}
		claims := token.Claims.(*auth.UserClaims)
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)
		h(w, r)
	}
}
func getClaimsFromContext(r *http.Request) (*auth.UserClaims, bool) {
	claims, ok := r.Context().Value("claims").(*auth.UserClaims)
	return claims, ok
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
		claims, ok := getClaimsFromContext(r)
		if !ok {
			jsonMessage(w, http.StatusInternalServerError, "Could not get claims from context")
			return
		}
		if !claims.Admin {
			jsonMessage(w, http.StatusForbidden, "You are not allowed to access this resource")
			return
		}
		s.handleGetAccounts(w, r)
	default:
		jsonMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}
func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsFromContext(r)
	if !ok {
		jsonMessage(w, http.StatusInternalServerError, "Could not get claims from context")
		return
	}
	id := mux.Vars(r)["id"]
	switch r.Method {
	case http.MethodGet:
		if claims.User.ID != id && !claims.Admin {
			jsonMessage(w, http.StatusUnauthorized, "You can only access your own account")
			return
		}
		s.handleGetAccountByID(w, r)
	case http.MethodPut:
		if claims.User.ID != id && !claims.Admin {
			jsonMessage(w, http.StatusUnauthorized, "You can only update your own account")
			return
		}
		s.handleUpdateAccount(w, r)
	case http.MethodDelete:
		if !claims.Admin {
			jsonMessage(w, http.StatusForbidden, "Only admins are allowed to delete accounts")
			return
		}
		s.handleDeleteAccount(w, r)
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
	_, err_ := auth.GenerateUserJWT(account)
	if err_ != nil {
		jsonMessage(w, http.StatusInternalServerError, err_.Error())
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
	//TODO: Implement transfer handler
	return
}

const welcomeMessage = `Welcome to the Bank JSON API Server! :)

Available endpoints:
GET /account - get all accounts
POST /account - create a new account
GET /account/{id} - get account by ID
PUT /account/{id} - update account by ID
DELETE /account/{id} - delete account by ID
POST /transfer - transfer money between accounts`
