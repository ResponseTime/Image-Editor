package main

import (
	"fmt"
	"main/Db"
	"main/Router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	Db.ConnectDb()
	Db.Init()
	router := Router.SetupRouter()
	router.Run(":8080")
}
