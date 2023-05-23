package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Tuzi07/solvify-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

type Database interface {
	GetProblem(id string) (models.Problem, error)
	AddProblem(problem models.Problem) error
	UpdateProblem(problem models.Problem) error
	DeleteProblem(id string) error
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

func (mgr *MongoDB) GetProblem(id string) (models.Problem, error) {
	collection := mgr.client.Database("solvify").Collection("problems")
	var config models.Problem
	err := collection.FindOne(context.Background(), bson.D{bson.E{Key: "id", Value: id}}).Decode(&config)
	if err != nil {
		return models.Problem{}, err
	}
	return config, nil
}

func (mgr *MongoDB) AddProblem(problem models.Problem) error {
	collection := mgr.client.Database("solvify").Collection("problems")
	_, err := collection.InsertOne(context.Background(), problem)
	return err
}

func (mgr *MongoDB) UpdateProblem(problem models.Problem) error {
	collection := mgr.client.Database("solvify").Collection("problems")
	filter := bson.D{bson.E{Key: "id", Value: problem.ID}}
	update := bson.D{bson.E{Key: "$set", Value: problem}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (mgr *MongoDB) DeleteProblem(id string) error {
	collection := mgr.client.Database("solvify").Collection("problems")
	_, err := collection.DeleteOne(context.Background(), bson.D{bson.E{Key: "id", Value: id}})
	return err
}
