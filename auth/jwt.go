package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mathis-k/bank-api/models"
	"log"
	"os"
	"time"
)

var mySigningKey = os.Getenv("JWT_SECRET")

const (
	INVALID_TOKEN  = "Invalid token"
	TOKEN_EXPIRED  = "Token has expired"
	INVALID_CLAIMS = "Invalid token claims"
	TOKEN_PARSE    = "Error parsing token"
)

type user struct {
	ID            string `json:"id"`
	AccountNumber uint64 `json:"account_number"`
}
type UserClaims struct {
	ISS   string `json:"iss"`
	Admin bool   `json:"admin"`
	User  user   `json:"user"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
}

func (u UserClaims) GetIssuer() (string, error) {
	if u.ISS == "" {
		return "", fmt.Errorf("No issuer set")
	}
	return u.ISS, nil
}

func (u UserClaims) IsAdmin() (bool, error) {
	return u.Admin, nil
}

func (u UserClaims) GetUser() (string, uint64, error) {
	if u.User == (user{}) {
		return "", 0, fmt.Errorf("No user set")
	}
	if u.User.ID == "" {
		return "", 0, fmt.Errorf("No user ID set")
	} else if u.User.AccountNumber == 0 {
		return "", 0, fmt.Errorf("No account number set")
	}
	return u.User.ID, u.User.AccountNumber, nil
}

func (u UserClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if u.Exp == 0 {
		return nil, fmt.Errorf("No expiration time set")
	}
	expirationTime := jwt.NewNumericDate(time.Unix(u.Exp, 0))
	return expirationTime, nil
}

func (u UserClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	if u.Iat == 0 {
		return nil, fmt.Errorf("No issued at time set")
	}
	issuedAt := jwt.NewNumericDate(time.Unix(u.Iat, 0))
	return issuedAt, nil
}

func (u UserClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (u UserClaims) GetSubject() (string, error) {
	return "", nil
}

func (u UserClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func GenerateUserJWT(account *models.Account) (string, error) {
	claims := UserClaims{
		ISS:   "bank-api",
		Admin: false,
		User: user{
			ID:            account.ID.Hex(),
			AccountNumber: account.AccountNumber,
		},
		Exp: time.Now().Add(time.Hour * 24).Unix(),
		Iat: time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	log.Printf("✔ JWT token created for account %d / user %v: %v", account.AccountNumber, account.ID.Hex(), signedToken)
	return signedToken, nil
}

func GenerateAdminJWT(account *models.Account) (string, error) {
	//TODO: Later
	return "", nil
}

func VerifyJWT(signedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims.Exp < time.Now().Unix() {
			return nil, fmt.Errorf(TOKEN_EXPIRED)
		}
		return token, nil
	} else {
		if !token.Valid {
			return nil, fmt.Errorf(INVALID_TOKEN)
		}
		return nil, fmt.Errorf(INVALID_CLAIMS)
	}
}
