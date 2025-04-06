package main

import (
	"crop-tracker/db"
	"crop-tracker/routers"
	"log"
)

func main() {
	database := db.InitDB("data.db")
	r := routers.SetupRouter(database)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
