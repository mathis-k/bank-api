package models

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type DB struct {
	Client *mongo.Client
	Db     *mongo.Database
}

const (
	ConnectionWarningTimeOut = 2 * time.Second
	CloseTimeOut             = 5 * time.Second
	CheckConnectionTimeOut   = 2 * time.Second
)

func (d *DB) Connect() error {
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

	d.Client = client
	d.Db = client.Database(database)
	log.Println("✔ Successfully Connected to MongoDB")

	return nil
}

func (d *DB) IsConnected() bool {
	if d.Db == nil || d.Client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), CheckConnectionTimeOut)
	defer cancel()

	if err := d.Client.Ping(ctx, nil); err != nil {
		return false
	} else {
		return true
	}
}

func (d *DB) Disconnect() error {
	if d.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), CloseTimeOut)
		defer cancel()
		if err := d.Client.Disconnect(ctx); err != nil {
			log.Println("✖ Error disconnecting from MongoDB")
			return err
		} else {
			log.Println("✔ Successfully disconnected from MongoDB")
			return nil
		}
	} else {
		log.Println("✖ MongoDB connection is not active")
		return utils.DATABASE_NOT_ACTIVVE
	}
}
