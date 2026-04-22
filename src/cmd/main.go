package main

import (
	"kasoka/src/config"
	"kasoka/src/db"
	"kasoka/src/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../config/.env")

	mongo := config.LoadMongo()

	client, err := db.ConnectMongo(mongo.Uri)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(mongo.DbName)
	log.Println("mongo connected", db.Name())


	server.ServerInt()
}