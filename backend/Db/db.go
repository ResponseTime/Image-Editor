package Db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	err    error
)

func ConnectDb() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
}

func GetClient() *mongo.Client {
	if err != nil {
		return nil
	}
	return client
}
