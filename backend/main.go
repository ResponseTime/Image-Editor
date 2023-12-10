package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)
var client *mongo.Client
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	connectDb()
	
	router := setupRouter()
	
	router.Run(":8080")
}