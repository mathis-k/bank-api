package models

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/mathis-k/bank-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Password  string `bson:"password" json:"password" validate:"required,min=6"`
}

type UserUpdate struct {
	FirstName string `bson:"first_name" json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `bson:"last_name" json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     string `bson:"email" json:"email" validate:"omitempty,email"`
}

func ValidateUserRequest(request *UserRequest) error {
	validate := validator.New()
	return validate.Struct(request)
}

func ValidateUserUpdate(request *UserUpdate) error {
	validate := validator.New()
	return validate.Struct(request)
}

func (db *DB) CreateUser(userRequest *UserRequest) (*User, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
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
		return nil, err
	}
	return user, nil
}

func (db *DB) GetUserById(id primitive.ObjectID) (*User, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
	user := &User{}
	err := db.Db.Collection("users").FindOne(context.TODO(), primitive.M{"_id": id}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (*User, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
	user := &User{}
	err := db.Db.Collection("users").FindOne(context.TODO(), primitive.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) AddAccountToUser(user *User, aId primitive.ObjectID) error {
	if !db.IsConnected() {
		return utils.DATABASE_NOT_ACTIVVE
	}
	user.Accounts = append(user.Accounts, aId)
	_, err := db.Db.Collection("users").UpdateOne(context.TODO(), primitive.M{"_id": user.ID}, primitive.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateUser(user *User, userUpdate *UserUpdate) error {
	if !db.IsConnected() {
		return utils.DATABASE_NOT_ACTIVVE
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
	_, err := db.Db.Collection("users").UpdateOne(context.TODO(), primitive.M{"_id": user.ID}, primitive.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteUser(uId primitive.ObjectID) error {
	if !db.IsConnected() {
		return utils.DATABASE_NOT_ACTIVVE
	}
	_, err := db.Db.Collection("users").DeleteOne(context.TODO(), primitive.M{"_id": uId})
	if err != nil {
		return err
	}
	return nil
}
