package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var DB *mongo.Database

func ConnectDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	DB = client.Database("auth_service") // Select database
}
