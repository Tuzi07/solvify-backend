package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

type Database interface {
	ProblemDatabase
}

func NewMongoDB() (*MongoDB, error) {
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New("MONGODB_URI environment variable not set")
	}
	connectionOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, connectionOptions)
	if err != nil {
		log.Printf("could not connect to mongodb: %v\n", err)
	}

	mongoDB := &MongoDB{client: client}
	return mongoDB, err
}
