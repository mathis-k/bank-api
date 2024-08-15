package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mathis-k/bank-api/models"
	"log"
	"os"
)

var mySigningKey = os.Getenv("JWT_SECRET")

func GenerateUserJWT(account *models.Account) (string, error) {
	claims := jwt.MapClaims{
		"admin":          false,
		"id":             account.ID,
		"account_number": account.AccountNumber,
		"balance":        account.Balance,
		"first_name":     account.FirstName,
		"last_name":      account.LastName,
		"email":          account.Email,
		"created_at":     account.CreatedAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateAdminJWT(account *models.Account) (string, error) {
	claims := jwt.MapClaims{
		"admin":          true,
		"id":             account.ID,
		"account_number": account.AccountNumber,
		"balance":        account.Balance,
		"first_name":     account.FirstName,
		"last_name":      account.LastName,
		"email":          account.Email,
		"created_at":     account.CreatedAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyJWT(signedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if token.Valid {
		log.Println("Token is valid")
	} else {
		log.Println("Token is invalid")
	}
	return token, nil
}
