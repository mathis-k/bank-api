package models

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mathis-k/bank-api/middleware"
	"github.com/mathis-k/bank-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
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

func (db *DB) Connect() error {
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

	db.Client = client
	db.Db = client.Database(database)
	log.Println("✔ Successfully Connected to MongoDB")

	return nil
}

func (db *DB) IsConnected() bool {
	if db.Db == nil || db.Client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), CheckConnectionTimeOut)
	defer cancel()

	if err := db.Client.Ping(ctx, nil); err != nil {
		return false
	} else {
		return true
	}
}

func (db *DB) Disconnect() error {
	if db.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), CloseTimeOut)
		defer cancel()
		if err := db.Client.Disconnect(ctx); err != nil {
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

func (db *DB) CheckAccountPermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaimsFromContext(r)
		if !ok {
			utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
			return
		}

		user, err := db.GetUserById(claims.User_Id)
		if err != nil {
			utils.ErrorMessage(w, http.StatusPreconditionFailed, err)
		}

		vars := mux.Vars(r)
		accountNumber_str, ok := vars["number"]
		if !ok {
			utils.ErrorMessage(w, http.StatusBadRequest, utils.MISSING_ACCOUNT_NUMBER)
			return
		}
		accountNumber, err := utils.StringToUint64(accountNumber_str)
		if err != nil {
			utils.ErrorMessage(w, http.StatusBadRequest, err)
			return
		}

		account, err := db.GetAccountByAccountNumber(accountNumber)
		if err != nil {
			utils.ErrorMessage(w, http.StatusForbidden, err)
			return
		}
		if !user.HasAccount(account.ID) {
			utils.ErrorMessage(w, http.StatusNotFound, utils.ACCOUNT_NOT_FOUND)
			return
		}
		ctx := context.WithValue(r.Context(), "account", account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
