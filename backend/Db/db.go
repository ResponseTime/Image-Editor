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
	Client *mongo.Client
)

func ConnectDb() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
}

func Init() {
	if err != nil {
		Client = nil
	}
	Client = client
}
