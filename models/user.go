package models

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mathis-k/bank-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string               `bson:"first_name" json:"first_name"`
	LastName  string               `bson:"last_name" json:"last_name"`
	Email     string               `bson:"email" json:"email"`
	Password  string               `bson:"password" json:"password"`
	Accounts  []primitive.ObjectID `bson:"accounts" json:"accounts"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}
type UserRequest struct {
	FirstName string `bson:"first_name" json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `bson:"last_name" json:"last_name" validate:"required,min=2,max=50"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	Password  string `bson:"password" json:"password" validate:"required,min=8"`
}
type UserUpdate struct {
	FirstName string `bson:"first_name" json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `bson:"last_name" json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     string `bson:"email" json:"email" validate:"omitempty,email"`
}
type UserLogin struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required,min=8"`
}

func ValidateUserRequest(request *UserRequest) error {
	validate := validator.New()
	return validate.Struct(request)
}
func ValidateUserUpdate(request *UserUpdate) error {
	validate := validator.New()
	return validate.Struct(request)
}
func ValidateUserLogin(request *UserLogin) error {
	validate := validator.New()
	return validate.Struct(request)
}

func (u User) HasAccount(aId primitive.ObjectID) bool {
	for _, account := range u.Accounts {
		if account == aId {
			return true
		}
	}
	return false
}
func (db *DB) AddAccountToUser(uId primitive.ObjectID, aId primitive.ObjectID) error {
	user, err := db.GetUserById(uId)
	if err != nil {
		return err
	}
	user.Accounts = append(user.Accounts, aId)
	_, err = db.Db.Collection("users").UpdateOne(context.TODO(), primitive.M{"_id": user.ID}, primitive.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) RemoveAccountFromUser(uId primitive.ObjectID, aId primitive.ObjectID) error {
	user, err := db.GetUserById(uId)
	if err != nil {
		return err
	}
	for i, account := range user.Accounts {
		if account == aId {
			user.Accounts = append(user.Accounts[:i], user.Accounts[i+1:]...)
			break
		}
	}
	_, err = db.Db.Collection("users").UpdateOne(context.TODO(), primitive.M{"_id": user.ID}, primitive.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) CreateUser(userRequest *UserRequest) (*User, error) {
	password, err := utils.HashPassword(userRequest.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:        primitive.NewObjectID(),
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  password,
		Accounts:  []primitive.ObjectID{},
		CreatedAt: time.Now(),
	}
	_, err = db.Db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, utils.EMAIL_ALREADY_EXISTS
		}
		return nil, err
	}
	return user, nil
}
func (db *DB) GetUserById(id primitive.ObjectID) (*User, error) {
	user := &User{}
	err := db.Db.Collection("users").FindOne(context.TODO(), primitive.M{"_id": id}).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.USER_NOT_FOUND
		}
		return nil, err
	}
	return user, nil
}
func (db *DB) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.Db.Collection("users").FindOne(context.TODO(), primitive.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (db *DB) UpdateUser(uId primitive.ObjectID, userUpdate *UserUpdate) (*User, error) {
	user, err := db.GetUserById(uId)
	if err != nil {
		return nil, err
	}
	if userUpdate.FirstName != "" {
		user.FirstName = userUpdate.FirstName
	}
	if userUpdate.LastName != "" {
		user.LastName = userUpdate.LastName
	}
	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}
	updatedUser, err := db.Db.Collection("users").UpdateOne(context.TODO(), primitive.M{"_id": user.ID}, primitive.M{"$set": user})
	if err != nil {
		return nil, err
	}
	if updatedUser.MatchedCount == 0 {
		return nil, utils.USER_NOT_FOUND
	}
	return user, nil
}
func (db *DB) DeleteUser(uId primitive.ObjectID) error {
	_, err := db.Db.Collection("users").DeleteOne(context.TODO(), primitive.M{"_id": uId})
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) LoginUser(userLogin *UserLogin) (*User, error) {
	user, err := db.GetUserByEmail(userLogin.Email)
	if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(userLogin.Password, user.Password) {
		return nil, utils.INVALID_CREDENTIALS
	}
	return user, nil
}
