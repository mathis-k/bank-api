package db

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/models"
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

func (m *MongoDB) Connect() error {
	if err := godotenv.Load(); err != nil {
		log.Println("✖ No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Println("✖ Error connecting to MongoDB")
		return err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("✖ Error disconnecting from MongoDB")
			panic(err)
		}
	}()
	m.Client = client
	m.db = client.Database("TestDB")
	log.Println("✔ Sucessfully Connected to MongoDB")
	return nil
}

func (m *MongoDB) isConnected() bool {
	if m.db == nil || m.Client == nil {
		log.Println("⚠ MongoDB connection is not established")
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := m.Client.Ping(ctx, nil); err != nil {
		log.Println("⚠ MongoDB connection is not established")
		return false
	} else {
		log.Println("ℹ MongoDB connection is established")
		return true
	}
}

func (m *MongoDB) CreateAccount(a *models.Account) error {
	return nil
}

func (m *MongoDB) GetAccounts() ([]*models.Account, error) {
	return nil, nil
}

func (m *MongoDB) GetAccountByID(id int) (*models.Account, error) {
	return nil, nil
}

func (m *MongoDB) DeleteAccount(id int) error {
	return nil
}

func (m *MongoDB) UpdateAccount(id int, a *models.Account) error {
	return nil
}
