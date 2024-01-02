package main

import (
	"fmt"
	"main/Db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	Db.ConnectDb()

	router := setupRouter()

	router.Run(":8080")
}
