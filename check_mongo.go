package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb+srv://leykun:QO0G695rmMbRSrD2@cluster0.f76x9.mongodb.net/?appName=Cluster0"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	dbName := "ethiopian_financial_news"
	db := client.Database(dbName)

	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Collections in %s:\n", dbName)
	for _, coll := range collections {
		fmt.Printf(" - %s\n", coll)

		// count documents
		count, _ := db.Collection(coll).CountDocuments(ctx, bson.D{})
		fmt.Printf("   count: %d\n", count)
	}
}
