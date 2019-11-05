package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbURL string
var dbName string
var client *mongo.Client
var db *mongo.Database

func init() {
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		dbURL = os.Getenv("DB_URL")
		dbName = os.Getenv("DB_NAME")
	}

	clientOptions := options.Client().ApplyURI(dbURL)

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(dbName)
}

func main() {
	go spinner(100 * time.Millisecond)

	// Check the connection
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\rSuccessfully connected to MongoDB.\n")
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
