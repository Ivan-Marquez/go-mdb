package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setIndexModel(keys primitive.M, opts *options.IndexOptions) mongo.IndexModel {
	index := mongo.IndexModel{}

	index.Keys = keys

	if opts != nil {
		index.Options = opts
	}

	return index
}

// user indexes
func getUserIndexes() []mongo.IndexModel {
	email := setIndexModel(bson.M{
		"email": 1,
	}, nil)

	searchIndexName := "searchIndex"
	searchIndex := setIndexModel(bson.M{
		"lastName":   "text",
		"firstName":  "text",
		"middleName": "text",
	}, &options.IndexOptions{
		Name: &searchIndexName,
		Weights: bson.M{
			"lastName":   16,
			"firstName":  16,
			"middleName": 8,
		},
	})

	return []mongo.IndexModel{
		email,
		searchIndex,
	}
}

// activity indexes
func getActivityIndexes() []mongo.IndexModel {
	locationName := setIndexModel(bson.M{
		"locationName": 1,
	}, nil)

	distance := setIndexModel(bson.M{
		"distance": 1,
	}, nil)

	activityID := setIndexModel(bson.M{
		"activityId": 1,
	}, nil)

	return []mongo.IndexModel{
		locationName,
		distance,
		activityID,
	}
}
