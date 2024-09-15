package utils

import (
	"fmt"
	"github.com/mathis-k/bank-api/middleware"
	"net/http"
	"time"
)

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%02dh %02dmin", hours, minutes)
}

func GetClaimsFromContext(r *http.Request) (*middleware.UserClaims, bool) {
	claims, ok := r.Context().Value("claims").(*middleware.UserClaims)
	return claims, ok
}
