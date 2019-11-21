package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	*mongo.Database
}

func (db *database) setIndexes(collection string, indexes []mongo.IndexModel, c chan<- string) {
	col := db.Collection(collection)
	col.Indexes().CreateMany(context.Background(), indexes)

	c <- fmt.Sprintf("\rCreated indexes for %s\n", collection)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connect(dbURL, dbName string) (client *mongo.Client, err error) {
	clientOptions := options.Client().ApplyURI(dbURL)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	client, err := connect(dbURL, dbName)
	if err != nil {
		log.Fatal("Error connecting to MongoDB instance:", err)
	}

	defer client.Disconnect(context.TODO())

	db := new(database)
	db.Database = client.Database(dbName)

	collections := map[string][]mongo.IndexModel{
		"user":     getUserIndexes(),
		"activity": getActivityIndexes(),
	}

	fmt.Printf("\rGenerating indexes:\n")
	ch := make(chan string, len(collections))

	for collection, indexes := range collections {
		go db.setIndexes(collection, indexes, ch)
	}

	for i := 0; i < len(collections); i++ {
		fmt.Printf(<-ch)
	}
}
