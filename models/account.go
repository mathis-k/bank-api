package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"regexp"
	"time"
	"unicode"
)

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccountNumber uint64             `bson:"account_number" json:"account_number"`
	Balance       float64            `bson:"balance" json:"balance"`
	FirstName     string             `bson:"first_name" json:"first_name"`
	LastName      string             `bson:"last_name" json:"last_name"`
	Email         string             `bson:"email" json:"email"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type AccountRequest struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
}

const MaxNameLength = 50
const MinNameLength = 2
const MaxEmailLength = 100
const MinEmailLength = 6
const InvalidFirstName = "invalid first name"
const InvalidLastName = "invalid last name"
const InvalidEmail = "invalid email"

func NewAccount(req *AccountRequest) (*Account, error) {
	if err := IsValidAccountRequest(req); err != nil {
		return nil, err
	}
	return &Account{
		ID:            primitive.NewObjectID(),
		AccountNumber: GenerateUniqueAccountNumber(),
		Balance:       0.0,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		CreatedAt:     time.Now().Local(),
	}, nil
}
func IsValidAccountRequest(req *AccountRequest) error {
	if !isValidName(req.FirstName) || !isValidLength(req.FirstName, MinNameLength, MaxNameLength) {
		return fmt.Errorf(InvalidFirstName)
	} else if !isValidName(req.LastName) || !isValidLength(req.LastName, MinNameLength, MaxNameLength) {
		return fmt.Errorf(InvalidLastName)
	} else if !isValidEmail(req.Email) || !isValidLength(req.Email, MinEmailLength, MaxEmailLength) {
		return fmt.Errorf(InvalidEmail)
	}
	return nil
}

func GenerateUniqueAccountNumber() uint64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint64()
}

func isValidName(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidLength(value string, min, max int) bool {
	length := len(value)
	return length >= min && length <= max
}
