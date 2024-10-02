package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccountNumber uint64             `bson:"account_number" json:"account_number"`
	Balance       float64            `bson:"balance" json:"balance"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

func (db *DB) CreateAccount() (*Account, error) {
	account := &Account{
		ID:            primitive.NewObjectID(),
		AccountNumber: uint64(time.Now().Unix()),
		Balance:       0.0,
		CreatedAt:     time.Now(),
	}
	_, err := db.Db.Collection("accounts").InsertOne(context.TODO(), account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (db *DB) GetAccountById(aId primitive.ObjectID) (*Account, error) {
	account := &Account{}
	err := db.Db.Collection("accounts").FindOne(context.TODO(), primitive.M{"_id": aId}).Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (db *DB) GetAccountByAccountNumber(accountNumber uint64) (*Account, error) {
	account := &Account{}
	err := db.Db.Collection("accounts").FindOne(context.TODO(), primitive.M{"account_number": accountNumber}).Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (db *DB) GetAccountsFromUser(uId primitive.ObjectID) ([]*Account, error) {
	user, err := db.GetUserById(uId)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for _, accountID := range user.Accounts {
		account := &Account{}
		err := db.Db.Collection("accounts").FindOne(context.TODO(), primitive.M{"_id": accountID}).Decode(account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
func (db *DB) DeleteAccount(aId primitive.ObjectID) error {
	_, err := db.Db.Collection("accounts").DeleteOne(context.TODO(), primitive.M{"_id": aId})
	if err != nil {
		return err
	}
	return nil
}
