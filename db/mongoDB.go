package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Database interface {
	Connect() error
	Disconnect() error
}

type MongoDB struct {
	Client *mongo.Client
	db     *mongo.Database
}

const (
	ConnectionWarningTimeOut = 2 * time.Second
	CloseTimeOut             = 5 * time.Second
	CheckConnectionTimeOut   = 2 * time.Second

	DataBaseNotActive = "MongoDB connection is not active"
)

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

func (m *MongoDB) Disconnect() error {
	if m.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), CloseTimeOut)
		defer cancel()
		if err := m.Client.Disconnect(ctx); err != nil {
			log.Println("✖ Error disconnecting from MongoDB")
			return err
		} else {
			log.Println("✔ Successfully disconnected from MongoDB")
			return nil
		}
	} else {
		log.Println("✖ MongoDB connection is not active")
		return fmt.Errorf(DataBaseNotActive)
	}
}
