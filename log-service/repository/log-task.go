package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Fungsi untuk menambahkan validator ke tasks-log
func addTaskCollectionValidator(ctx context.Context, client *mongo.Client, dbName, collectionName string) error {
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
			"required": []string{"userId", "taskId", "action", "timestamp"},
			"properties": bson.M{
				"userId": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"taskId": bson.M{
					"bsonType":    "long",
					"description": "must be an integer and is required",
				},
				"action": bson.M{
					"bsonType":    "string",
					"description": "must be a valid task action and is required",
				},
				"timestamp": bson.M{
					"bsonType":    "string",
					"description": "must be a string representing the timestamp and is required",
				},
			},
		},
	}

	// Gunakan perintah collMod untuk menambahkan validator pada koleksi yang sudah ada
	command := bson.D{
		{Key: keyCommandCollection, Value: collectionName},
		{Key: "validator", Value: validator},
	}

	// Jalankan perintah collMod
	return client.Database(dbName).RunCommand(ctx, command).Err()
}
