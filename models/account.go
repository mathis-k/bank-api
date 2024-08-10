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
	AccountNumber int64              `bson:"account_number" json:"account_number"`
	Balance       float64            `bson:"balance" json:"balance"`
	FirstName     string             `bson:"first_name" json:"first_name"`
	LastName      string             `bson:"last_name" json:"last_name"`
	Email         string             `bson:"email" json:"email"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type CreateAccountRequest struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
}

const MaxNameLength = 50
const MinNameLength = 2
const MaxEmailLength = 100
const MinEmailLength = 6

func NewAccount(req *CreateAccountRequest) (*Account, error) {
	if !isValidName(req.FirstName) || !isValidLength(req.FirstName, MinNameLength, MaxNameLength) {
		return nil, fmt.Errorf("invalid first name")
	} else if !isValidName(req.LastName) || !isValidLength(req.LastName, MinNameLength, MaxNameLength) {
		return nil, fmt.Errorf("invalid last name")
	} else if !isValidEmail(req.Email) || !isValidLength(req.Email, MinEmailLength, MaxEmailLength) {
		return nil, fmt.Errorf("invalid email")
	}
	return &Account{
		ID:            primitive.NewObjectID(),
		AccountNumber: GenerateUniqueAccountNumber(),
		Balance:       0.0,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		CreatedAt:     time.Now(),
	}, nil
}
func UpdateAccount(id string, req *CreateAccountRequest) (*Account, error) {
	if !isValidName(req.FirstName) || !isValidLength(req.FirstName, MinNameLength, MaxNameLength) {
		return nil, fmt.Errorf("invalid first name")
	} else if !isValidName(req.LastName) || !isValidLength(req.LastName, MinNameLength, MaxNameLength) {
		return nil, fmt.Errorf("invalid last name")
	} else if !isValidEmail(req.Email) || !isValidLength(req.Email, MinEmailLength, MaxEmailLength) {
		return nil, fmt.Errorf("invalid email")
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}
	return &Account{
		ID:        _id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}, nil
}
func GenerateUniqueAccountNumber() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(10000000000))
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
