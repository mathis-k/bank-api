package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/models"
	"go.mongodb.org/mongo-driver/bson"
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

const MaxAttempts = 3 /* Maximum number of attempts to generate a unique account number */
const ConnectionWarningTimeOut = 2 * time.Second

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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := m.Client.Ping(ctx, nil); err != nil {
		return false
	} else {
		return true
	}
}

func (m *MongoDB) CreateAccount(a *models.Account) error {
	if !m.isConnected() {
		return fmt.Errorf("MongoDB connection is not active")
	}
	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingAccount models.Account
	err := collection.FindOne(ctx, bson.M{"email": a.Email}).Decode(&existingAccount)
	if err == nil {
		return fmt.Errorf("an account with the email %s already exists", a.Email)
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	for i := 0; i < MaxAttempts; i++ {
		a.AccountNumber = models.GenerateUniqueAccountNumber()

		var existingAccountByNumber models.Account
		err := collection.FindOne(ctx, bson.M{"account_number": a.AccountNumber}).Decode(&existingAccountByNumber)
		if err == mongo.ErrNoDocuments {
			break
		}
		if i == MaxAttempts-1 {
			return fmt.Errorf("could not generate a unique account number after %d attempts, please try again", MaxAttempts)
		}
	}

	_, err2 := collection.InsertOne(ctx, a)
	if err2 != nil {
		return err2
	}
	return nil
}

func (m *MongoDB) GetAllAccounts() ([]*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf("MongoDB connection is not active")
	}
	collection := m.db.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("⚠ Getting All Accounts From MongoDB: Error closing cursor")
		}
	}(cursor, ctx)

	var accounts []*models.Account
	for cursor.Next(ctx) {
		var account models.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)

	}
	return accounts, nil
}

func (m *MongoDB) GetAccountByID(id int) (*models.Account, error) {
	if !m.isConnected() {
		return nil, fmt.Errorf("MongoDB connection is not active")
	}
	return nil, nil
}

func (m *MongoDB) DeleteAccount(id int) error {
	if !m.isConnected() {
		return fmt.Errorf("MongoDB connection is not active")
	}
	return nil
}

func (m *MongoDB) UpdateAccount(id int, a *models.Account) error {
	if !m.isConnected() {
		return fmt.Errorf("MongoDB connection is not active")
	}
	return nil
}
