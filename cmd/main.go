package main

import (
	"log"

	"github.com/Tuzi07/solvify-backend/internal/api"
	"github.com/Tuzi07/solvify-backend/internal/db"
)

func main() {
	db, err := db.NewMongoDB()
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	server := api.NewServer(db)
	err = server.Start()
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
