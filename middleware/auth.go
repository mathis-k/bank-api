package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mathis-k/bank-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var mySigningKey = os.Getenv("JWT_SECRET")

const (
	EXPIRATION_TIME_USER = time.Hour * 24
)

type UserClaims struct {
	User_Id primitive.ObjectID `json:"user"`
	Valid   bool               `json:"valid"`
	Exp     int64              `json:"exp"`
	Iat     int64              `json:"iat"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorMessage(w, http.StatusUnauthorized, utils.MISSING_AUTH_HEADER)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		token, err := VerifyJWT(tokenString)
		if err != nil {
			if errors.Is(err, utils.TOKEN_EXPIRED) || errors.Is(err, utils.INVALID_TOKEN) {
				utils.ErrorMessage(w, http.StatusUnauthorized, err)
			} else {
				utils.ErrorMessage(w, http.StatusBadRequest, err)
			}
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateUserJWT(uId primitive.ObjectID) (string, error) {
	claims := UserClaims{
		User_Id: uId,
		Valid:   true,
		Exp:     time.Now().Add(EXPIRATION_TIME_USER).Unix(),
		Iat:     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	log.Printf("ℹ New JWT token created for user %v (Valid for %s): %v", claims.User_Id, utils.FormatDuration(EXPIRATION_TIME_USER), signedToken)
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
			return nil, utils.TOKEN_EXPIRED
		} else if !claims.Valid {
			return nil, utils.INVALID_TOKEN
		}
		return token, nil
	} else {
		if !token.Valid {
			return nil, utils.INVALID_TOKEN
		}
		return nil, utils.INVALID_CLAIMS
	}
}

func GetClaimsFromContext(r *http.Request) (*UserClaims, bool) {
	claims, ok := r.Context().Value("claims").(*UserClaims)
	return claims, ok
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
