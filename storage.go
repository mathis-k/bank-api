package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Storage interface {
	CreateAccount(a *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
	DeleteAccount(id int) error
	UpdateAccount(id int, a *Account) error
}

type MongoDBStorage struct {
	client *mongo.Client
	db     *mongo.Database
}

func (m *MongoDBStorage) Connect() error {
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
	m.client = client
	m.db = client.Database("TestDB")
	log.Println("✔ Sucessfully Connected to MongoDB")
	return nil
}

func (m *MongoDBStorage) isConnected() bool {
	if m.db == nil || m.client == nil {
		log.Println("⚠ MongoDB connection is not established")
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := m.client.Ping(ctx, nil); err != nil {
		log.Println("⚠ MongoDB connection is not established")
		return false
	} else {
		log.Println("ℹ MongoDB connection is established")
		return true
	}
}
