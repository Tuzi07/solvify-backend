package main

import (
	"log"

	"github.com/Tuzi07/solvify-backend/internal/database"
	"github.com/Tuzi07/solvify-backend/internal/server"
)

func main() {
	db, err := database.NewMongoDB()
	if err != nil {
		log.Panic(err)
	}
	err = server.Run(db)
	log.Panic(err)
}
