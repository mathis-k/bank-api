package models

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
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

type TransactionRequest struct {
	Type        TransactionType    `bson:"type" json:"type" validate:"required"`
	Amount      float64            `bson:"amount" json:"amount" validate:"required,gt=0,lte=10000"`
	FromAccount primitive.ObjectID `bson:"from_account" json:"from_account"`
	ToAccount   uint64             `bson:"to_account" json:"to_account"`
	ToAccountID primitive.ObjectID
}

func ValidateTransactionRequest(request *TransactionRequest) error {
	validate := validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		req := sl.Current().Interface().(TransactionRequest)

		switch req.Type {
		case Transfer:
			if req.FromAccount == primitive.NilObjectID {
				sl.ReportError(req.FromAccount, "FromAccount", "from_account", "requiredForTransfer", "")
			}
			if req.ToAccountID == primitive.NilObjectID {
				sl.ReportError(req.ToAccountID, "ToAccountID", "to_account_id", "requiredForTransfer", "")
			}
		case Deposit:
			if req.ToAccountID == primitive.NilObjectID {
				sl.ReportError(req.ToAccountID, "ToAccountID", "to_account_id", "requiredForDeposit", "")
			}
			if req.FromAccount != primitive.NilObjectID {
				sl.ReportError(req.FromAccount, "FromAccount", "from_account", "shouldBeEmptyForDeposit", "")
			}
		case Payout:
			if req.FromAccount == primitive.NilObjectID {
				sl.ReportError(req.FromAccount, "FromAccount", "from_account", "requiredForPayout", "")
			}
			if req.ToAccountID != primitive.NilObjectID {
				sl.ReportError(req.ToAccountID, "ToAccountID", "to_account_id", "shouldBeEmptyForPayout", "")
			}
		}
	}, TransactionRequest{})

	return validate.Struct(request)
}

func (db *DB) CreateTransaction(transactionRequest *TransactionRequest) (*Transaction, error) {

	switch transactionRequest.Type {
	case Deposit:
		err := db.MakeDeposit(transactionRequest.Amount, transactionRequest.ToAccountID)
		if err != nil {
			return nil, err
		}
	case Payout:
		err := db.MakePayout(transactionRequest.Amount, transactionRequest.FromAccount)
		if err != nil {
			return nil, err
		}
	case Transfer:
		err := db.MakeTransfer(transactionRequest.Amount, transactionRequest.FromAccount, transactionRequest.ToAccountID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, utils.INVALID_TRANSACTION_TYPE
	}

	transaction := &Transaction{
		ID:          primitive.NewObjectID(),
		Type:        transactionRequest.Type,
		Amount:      transactionRequest.Amount,
		FromAccount: transactionRequest.FromAccount,
		ToAccount:   transactionRequest.ToAccountID,
		CreatedAt:   time.Now(),
	}
	_, err := db.Db.Collection("transactions").InsertOne(context.TODO(), transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
func (db *DB) MakeDeposit(amount float64, toAccount primitive.ObjectID) error {
	filter := primitive.M{"_id": toAccount}
	update := primitive.M{
		"$inc": primitive.M{"balance": amount},
	}

	_, err := db.Db.Collection("accounts").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) MakePayout(amount float64, fromAccount primitive.ObjectID) error {
	filter := primitive.M{
		"_id":     fromAccount,
		"balance": primitive.M{"$gte": amount},
	}

	update := primitive.M{
		"$inc": primitive.M{"balance": -amount},
	}

	account := db.Db.Collection("accounts").FindOneAndUpdate(context.TODO(), filter, update)
	if account.Err() != nil {
		if errors.Is(account.Err(), mongo.ErrNoDocuments) {
			return utils.INSUFFICIENT_FUNDS
		}
		return account.Err()
	}
	return nil
}
func (db *DB) MakeTransfer(amount float64, from_aId primitive.ObjectID, to_aId primitive.ObjectID) error {
	session, err := db.Db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		from_filter := primitive.M{
			"_id":     from_aId,
			"balance": primitive.M{"$gte": amount},
		}

		to_filter := primitive.M{"_id": to_aId}

		from_update := primitive.M{
			"$inc": primitive.M{"balance": -amount},
		}

		to_update := primitive.M{
			"$inc": primitive.M{"balance": amount},
		}

		fromAccount := db.Db.Collection("accounts").FindOneAndUpdate(sessCtx, from_filter, from_update)
		if fromAccount.Err() != nil {
			if errors.Is(fromAccount.Err(), mongo.ErrNoDocuments) {
				return nil, utils.INSUFFICIENT_FUNDS
			}
			return nil, fromAccount.Err()
		}

		_, err = db.Db.Collection("accounts").UpdateOne(sessCtx, to_filter, to_update)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(context.TODO(), callback)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) GetTransactionById(tId primitive.ObjectID) (*Transaction, error) {
	transaction := &Transaction{}
	err := db.Db.Collection("transactions").FindOne(context.TODO(), primitive.M{"_id": tId}).Decode(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
func (db *DB) GetTransactionsFromAccount(aId primitive.ObjectID) ([]*Transaction, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_account": aId},
			{"to_account": aId},
		},
	}
	cursor, err := db.Db.Collection("transactions").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
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
func (db *DB) GetTransactionsFromUser(uId primitive.ObjectID) ([]*Transaction, error) {
	user, err := db.GetUserById(uId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"$or": []bson.M{
			{"from_account": bson.M{"$in": user.Accounts}},
			{"to_account": bson.M{"$in": user.Accounts}},
		},
	}

	cursor, err := db.Db.Collection("transactions").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
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
