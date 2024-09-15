package utils

import (
	"fmt"
	"net/http"
)

var (
	DATABASE_NOT_ACTIVVE = fmt.Errorf("MongoDB connection is not active")
	INVALID_TOKEN        = fmt.Errorf("invalid token")
	TOKEN_EXPIRED        = fmt.Errorf("token has expired")
	INVALID_CLAIMS       = fmt.Errorf("invalid token claims")
	MISSING_AUTH_HEADER  = fmt.Errorf("Missing Authorization header")
)

func ErrorMessage(w http.ResponseWriter, code int, error error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, error.Error())))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
