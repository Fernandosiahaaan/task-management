package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func addUserActionValidator(ctx context.Context, client *mongo.Client, dbName, collectionName string) error {
	findCollections, err := client.Database(dbName).ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	keyCommandCollection := ""
	if len(findCollections) > 0 {
		keyCommandCollection = "collMod"
	} else {
		keyCommandCollection = "create"
	}

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"userId", "action", "timestamp"},
			"properties": bson.M{
				"userId": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"action": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"timestamp": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
			},
		},
	}

	command := bson.D{
		{Key: keyCommandCollection, Value: collectionName},
		{Key: "validator", Value: validator},
	}

	return client.Database(dbName).RunCommand(ctx, command).Err()
}
