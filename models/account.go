package models

import (
	"context"
	"github.com/mathis-k/bank-api/utils"
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
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
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

func (db *DB) GetAccountById(id primitive.ObjectID) (*Account, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
	account := &Account{}
	err := db.Db.Collection("accounts").FindOne(context.TODO(), primitive.M{"_id": id}).Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (db *DB) GetAccountByAccountNumber(accountNumber uint64) (*Account, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
	account := &Account{}
	err := db.Db.Collection("accounts").FindOne(context.TODO(), primitive.M{"account_number": accountNumber}).Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (db *DB) GetAccountsFromUser(user User) ([]*Account, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
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

func (db *DB) DeleteAccount(account *Account) error {
	if !db.IsConnected() {
		return utils.DATABASE_NOT_ACTIVVE
	}
	_, err := db.Db.Collection("accounts").DeleteOne(context.TODO(), primitive.M{"_id": account.ID})
	if err != nil {
		return err
	}
	return nil
}
