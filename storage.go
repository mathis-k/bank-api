package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Storage interface {
	CreateAccount(a *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
	DeleteAccount(id int) error
	UpdateAccount(id int, a *Account) error
}

type MongoDBStorage struct {
}

func (m *MongoDBStorage) Connect() (*mongo.Database, error) {
	//TODO
	if err := godotenv.Load(); err != nil {
		log.Println("✖ No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Println("✖ Error connecting to MongoDB")
		return nil, err
	} else {
		log.Println("✔ Sucessfully Connected to MongoDB")
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("✖ Error disconnecting from MongoDB")
			panic(err)
		}
	}()
	return client.Database("TestDB"), nil
}
