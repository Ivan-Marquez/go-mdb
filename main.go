package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	*mongo.Database
}

func (db *database) setIndexes(collection string, indexes []mongo.IndexModel) {
	col := db.Collection(collection)
	col.Indexes().CreateMany(context.Background(), indexes)
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
	var wg sync.WaitGroup
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

	wg.Add(len(collections))

	fmt.Printf("\rCPU's: %d\n", runtime.NumCPU())
	fmt.Printf("\rGenerating indexes:\n")

	for collection, indexes := range collections {
		go func(col string, idx []mongo.IndexModel) {
			defer wg.Done()

			db.setIndexes(col, idx)
			fmt.Printf("\rCreated indexes for %s\n", col)
		}(collection, indexes)
	}

	wg.Wait()
}
