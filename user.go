package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type credentialsValidation struct {
	EntryCount int `bson:"entryCount"`
}

type user struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	Email     string
	Phone     string
	CreatedOn primitive.DateTime `bson:"createdOn"`
}
