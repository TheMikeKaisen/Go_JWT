package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {

	// load env
	loadEnvErr := godotenv.Load(".env")
	if loadEnvErr != nil {
		fmt.Println("Error while loading env")
		return nil
	}

	mongoUri := os.Getenv("MONGO_URI")

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to mongo db
	client, mongoErr := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if mongoErr != nil {
		fmt.Println("Mongo connection Error")
		return nil
	}
	fmt.Println("Connected to MongoDB!")
	return client

}

var Client *mongo.Client = DBinstance()

func OpenCollection(collectionName string) *mongo.Collection {
	// create a collection
	collection := Client.Database("cluster0").Collection(collectionName)
	return collection
}
