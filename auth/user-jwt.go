package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mathis-k/bank-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
)

var mySigningKey = os.Getenv("JWT_SECRET")

const (
	INVALID_TOKEN        = "invalid token"
	TOKEN_EXPIRED        = "token has expired"
	INVALID_CLAIMS       = "invalid token claims"
	EXPIRATION_TIME_USER = time.Hour * 24
)

type UserClaims struct {
	User_Id primitive.ObjectID `json:"user"`
	Exp     int64              `json:"exp"`
	Iat     int64              `json:"iat"`
}

func GenerateUserJWT(user *models.User) (string, error) {
	claims := UserClaims{
		User_Id: user.ID,
		Exp:     time.Now().Add(EXPIRATION_TIME_USER).Unix(),
		Iat:     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	log.Printf("ℹ New JWT token created for user %v (Valid for %s): %v", user.ID, formatDuration(EXPIRATION_TIME_USER), signedToken)
	return signedToken, nil
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

func RefreshJWT(signedToken string) (string, error) {
	token, err := VerifyJWT(signedToken)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", fmt.Errorf(INVALID_CLAIMS)
	}
	claims.Exp = time.Now().Add(EXPIRATION_TIME_USER).Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = newToken.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	log.Printf("ℹ JWT token refreshed for user %v (Valid for %s): %v", claims.User_Id, formatDuration(EXPIRATION_TIME_USER), signedToken)
	return signedToken, nil
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%02dh %02dmin", hours, minutes)
}

func (u UserClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if u.Exp == 0 {
		return nil, fmt.Errorf("no expiration time set")
	}
	expirationTime := jwt.NewNumericDate(time.Unix(u.Exp, 0))
	return expirationTime, nil
}

func (u UserClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	if u.Iat == 0 {
		return nil, fmt.Errorf("no issued at time set")
	}
	issuedAt := jwt.NewNumericDate(time.Unix(u.Iat, 0))
	return issuedAt, nil
}

func (u UserClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (u UserClaims) GetIssuer() (string, error) {
	return "", nil
}

func (u UserClaims) GetSubject() (string, error) {
	return "", nil
}

func (u UserClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}
