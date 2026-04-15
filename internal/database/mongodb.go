package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"stocet-core-service/internal/config"
)

var MongoDBClient *mongo.Client

// InitMongoDB initializes the MongoDB connection for reading news/market data.
// On subsequent calls (e.g. Lambda warm starts) the existing connection is reused.
func InitMongoDB() error {
	if MongoDBClient != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDBURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	MongoDBClient = client
	log.Println("Connected to MongoDB")
	return nil
}

// CloseMongoDB closes the MongoDB connection
func CloseMongoDB() error {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return MongoDBClient.Disconnect(ctx)
	}
	return nil
}
