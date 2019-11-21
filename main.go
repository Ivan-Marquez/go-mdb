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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connect(dbURL string) (client *mongo.Client, err error) {
	clientOptions := options.Client().ApplyURI(dbURL)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	go spinner(100 * time.Millisecond)

	dbURL := os.Getenv("DB_URL")
	client, err := connect(dbURL)

	if err != nil {
		log.Fatal("Error connecting to MongoDB instance:", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

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
