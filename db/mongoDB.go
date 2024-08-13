package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type MongoDB struct {
	Client *mongo.Client
	db     *mongo.Database
}

const MongoDBCompassURI = "mongodb+srv://TestUser:TestPassword@testcluster.ytxup.mongodb.net/"

const MaxAttempts = 3 /* Maximum number of attempts to generate a unique account number */
const ConnectionWarningTimeOut = 2 * time.Second
const GetTimeOut = 5 * time.Second
const InsertTimeOut = 5 * time.Second
const CloseTimeOut = 5 * time.Second
const CheckConnectionTimeOut = 2 * time.Second
const DeleteTimeOut = 5 * time.Second

const DataBaseNotActive = "MongoDB connection is not active"
const InvalidID = "invalid id"
const InvalidAccountNumber = "invalid Account-Number"
const NoAccountFound = "no account found"

func (m *MongoDB) Connect() error {
	if err := godotenv.Load(); err != nil {
		log.Println("✖ No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	database := os.Getenv("MONGODB_DB")

	startTime := time.Now()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("✖ Error connecting to MongoDB")
		return err
	}

	elapsedTime := time.Since(startTime)

	if elapsedTime > ConnectionWarningTimeOut {
		log.Printf("⚠ Connection to MongoDB is taking longer than expected: %v", elapsedTime)
	}

	m.Client = client
	m.db = client.Database(database)
	log.Println("✔ Successfully Connected to MongoDB")

	return nil
}

func (m *MongoDB) isConnected() bool {
	if m.db == nil || m.Client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), CheckConnectionTimeOut)
	defer cancel()

	if err := m.Client.Ping(ctx, nil); err != nil {
		return false
	} else {
		return true
	}
}

func (m *MongoDB) CreateAccount(req *models.AccountRequest) (*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf(DataBaseNotActive)
	}

	a, err := models.NewAccount(req)
	if err != nil {
		return nil, err
	}

	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), InsertTimeOut)
	defer cancel()

	var existingAccount models.Account
	err = collection.FindOne(ctx, bson.M{"email": a.Email}).Decode(&existingAccount)
	if err == nil {
		return nil, fmt.Errorf("an account with the email %s already exists", a.Email)
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	for i := 0; i < MaxAttempts; i++ {
		a.AccountNumber = models.GenerateUniqueAccountNumber()

		var existingAccountByNumber models.Account
		err := collection.FindOne(ctx, bson.M{"account_number": a.AccountNumber}).Decode(&existingAccountByNumber)
		if err == mongo.ErrNoDocuments {
			break
		}
		if i == MaxAttempts-1 {
			return nil, fmt.Errorf("could not generate a unique account number after %d attempts, please try again", MaxAttempts)
		}
	}

	_, err = collection.InsertOne(ctx, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (m *MongoDB) GetAllAccounts(maxResults uint64) ([]*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf(DataBaseNotActive)
	}

	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), GetTimeOut)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding accounts: %v", err)
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			log.Println("⚠ Error closing cursor:", closeErr)
		}
	}()
	var accounts []*models.Account
	for cursor.Next(ctx) {
		var account models.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, fmt.Errorf("error decoding account: %v", err)
		}
		accounts = append(accounts, &account)
		if uint64(len(accounts)) >= maxResults {
			break
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %v", err)
	}

	return accounts, nil
}

func (m *MongoDB) GetAccountByID(id string) (*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf(DataBaseNotActive)
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf(InvalidID)
	}

	collection := m.db.Collection("accounts")

	ctx, cancel := context.WithTimeout(context.Background(), GetTimeOut)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{
		"_id": _id,
	})
	if err != nil {
		return nil, fmt.Errorf("error finding account: %v", err)
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			log.Println("⚠ Error closing cursor:", closeErr, "whilst trying to fetch account for id:", id, "from MongoDB")
		}
	}()

	var account models.Account
	if cursor.Next(ctx) {
		if err := cursor.Decode(&account); err != nil {
			return nil, fmt.Errorf("error decoding account: %v", err)
		}
	} else {
		return nil, fmt.Errorf(NoAccountFound)
	}
	return &account, nil
}

func (m *MongoDB) DeleteAccount(id string) (*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf(DataBaseNotActive)
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf(InvalidID)
	}

	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), DeleteTimeOut)
	defer cancel()

	var account models.Account

	err = collection.FindOneAndDelete(ctx, bson.M{"_id": _id}).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(NoAccountFound)
		} else {
			return nil, fmt.Errorf("error deleting account: %v", err)
		}
	}

	return &account, nil
}

func (m *MongoDB) UpdateAccount(id string, req *models.AccountRequest) (*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf(DataBaseNotActive)
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf(InvalidID)
	}

	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), InsertTimeOut+GetTimeOut)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&models.Account{})
	if err == nil {
		return nil, fmt.Errorf("an account with the email %s already exists, please choose another email", req.Email)
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	var existingAccount models.Account
	err = collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&existingAccount)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(NoAccountFound)
		}
		return nil, fmt.Errorf("error finding account: %v", err)
	}

	if req.FirstName != "" {
		existingAccount.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existingAccount.LastName = req.LastName
	}
	if req.Email != "" {

		existingAccount.Email = req.Email
	}

	_, err = collection.ReplaceOne(ctx, bson.M{"_id": _id}, existingAccount)
	if err != nil {
		return nil, fmt.Errorf("error updating account: %v", err)
	}
	return &existingAccount, nil
}

func (m *MongoDB) Deposit(id string, account int64, amount float64) error {
	if !m.isConnected() {
		return fmt.Errorf(DataBaseNotActive)
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf(InvalidID)
	}
	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), InsertTimeOut+GetTimeOut)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{
		"_id":            _id,
		"account_number": account,
	}).Decode(&models.Account{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf(NoAccountFound)
		}
		return fmt.Errorf("error finding account: %v", err)
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": _id, "account_number": account}, bson.M{"$inc": bson.M{"balance": amount}})
	if err != nil {
		return fmt.Errorf("error updating account: %v", err)
	}
	return nil
}

func (m *MongoDB) Close() {
	if m.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), CloseTimeOut)
		defer cancel()
		if err := m.Client.Disconnect(ctx); err != nil {
			log.Println("✖ Error disconnecting from MongoDB")
		} else {
			log.Println("✔ Successfully disconnected from MongoDB")
		}
	}
}
