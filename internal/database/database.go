package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

type Database interface {
	create(obj any, collectionName string) (id string, err error)
	get(obj any, id string, collectionName string) error
	update(obj any, id string, collectionName string) error
	delete(id string, collectionName string) error
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

func (dbManager *MongoDB) create(obj any, collectionName string) (id string, err error) {
	collection := dbManager.client.Database("solvify").Collection(collectionName)
	result, err := collection.InsertOne(context.Background(), obj)

	id = result.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}

func (dbManager *MongoDB) get(obj any, id string, collectionName string) error {
	collection := dbManager.client.Database("solvify").Collection(collectionName)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

func (dbManager *MongoDB) update(obj any, id string, collectionName string) error {
	collection := dbManager.client.Database("solvify").Collection(collectionName)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.D{{Key: "$set", Value: obj}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (dbManager *MongoDB) delete(id string, collectionName string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := dbManager.client.Database("solvify").Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: objectID}}
	_, err = collection.DeleteOne(context.Background(), filter)
	return err
}
