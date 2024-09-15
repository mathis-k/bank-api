package models

import (
	"context"
	"github.com/mathis-k/bank-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TransactionType string

const (
	Deposit  TransactionType = "Deposit"
	Payout   TransactionType = "Payout"
	Transfer TransactionType = "Transfer"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        TransactionType    `bson:"type" json:"type"`
	Amount      float64            `bson:"amount" json:"amount"`
	FromAccount primitive.ObjectID `bson:"from_account" json:"from_account"`
	ToAccount   primitive.ObjectID `bson:"to_account" json:"to_account"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

func (db *DB) CreateTransaction(transactionType TransactionType, amount float64, fromAccount primitive.ObjectID, toAccount primitive.ObjectID) (*Transaction, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
	transaction := &Transaction{
		ID:          primitive.NewObjectID(),
		Type:        transactionType,
		Amount:      amount,
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		CreatedAt:   time.Now(),
	}
	_, err := db.Db.Collection("transactions").InsertOne(context.TODO(), transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (db *DB) GetTransactionById(id primitive.ObjectID) (*Transaction, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}
	transaction := &Transaction{}
	err := db.Db.Collection("transactions").FindOne(context.TODO(), primitive.M{"_id": id}).Decode(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (db *DB) GetTransactionsFromAccount(account primitive.ObjectID) ([]*Transaction, error) {
	if !db.IsConnected() {
		return nil, utils.DATABASE_NOT_ACTIVVE
	}

	filter := bson.M{
		"$or": []bson.M{
			{"from_account": account},
			{"to_account": account},
		},
	}
	cursor, err := db.Db.Collection("transactions").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.TODO())

	var transactions []*Transaction
	for cursor.Next(context.Background()) {
		transaction := &Transaction{}
		err := cursor.Decode(transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
