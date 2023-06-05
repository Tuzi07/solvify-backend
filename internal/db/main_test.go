package db

import (
	"log"
	"os"
	"testing"

	"github.com/Tuzi07/solvify-backend/internal/api"
)

var testServer *Server
var testDB *db.MongoDB

func TestMain(m *testing.M) {
	testDB, err := db.NewMongoDB()
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	testServer = NewServer(testDB)

	os.Exit(m.Run())

}
