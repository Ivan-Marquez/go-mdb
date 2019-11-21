package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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
	dbName := os.Getenv("DB_NAME")
	client, err := connect(dbURL)
	if err != nil {
		log.Fatal("Error connecting to MongoDB instance:", err)
	}

	db := client.Database(dbName)
	defer client.Disconnect(context.TODO())

	// get users
	var users []*user
	cur, err := db.Collection("user").Find(context.TODO(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var doc user
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		// do any cleanup to the data here

		users = append(users, &doc)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	newUsers := make([]interface{}, len(users))
	for i := range users {
		newUsers[i] = users[i]
	}

	insertOp, err := db.Collection("user_qa").InsertMany(context.TODO(), newUsers)
	fmt.Printf("\rInserted documents: %v\n", insertOp.InsertedIDs)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
