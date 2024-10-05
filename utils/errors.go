package utils

import (
	"fmt"
	"net/http"
)

var (
	DATABASE_NOT_ACTIVVE     = fmt.Errorf("mongoDB connection is not active")
	INVALID_TOKEN            = fmt.Errorf("invalid token")
	TOKEN_EXPIRED            = fmt.Errorf("token has expired")
	INVALID_CLAIMS           = fmt.Errorf("invalid token claims")
	INVALID_CREDENTIALS      = fmt.Errorf("invalid credentials")
	MISSING_AUTH_HEADER      = fmt.Errorf("missing authorization header")
	EMAIL_ALREADY_EXISTS     = fmt.Errorf("email already exists")
	USER_NOT_FOUND           = fmt.Errorf("user not found")
	TRANSACTION_NOT_FOUND    = fmt.Errorf("transaction not found")
	ACCOUNT_NOT_FOUND        = fmt.Errorf("account not found")
	INVALID_TRANSACTION_TYPE = fmt.Errorf("invalid transaction type")
	INSUFFICIENT_FUNDS       = fmt.Errorf("insufficient funds")
)

func ErrorMessage(w http.ResponseWriter, code int, error error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, error.Error())))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
