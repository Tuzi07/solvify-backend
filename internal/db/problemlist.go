package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProblemListDatabase interface {
	CreateProblemList(arg CreateProblemListParams) (ProblemList, error)
	GetProblemList(id string) (ProblemList, error)
	ListProblemLists(arg ListProblemListsParams) ([]ProblemList, error)
	UpdateProblemList(arg UpdateProblemListParams, id string) error
	DeleteProblemList(id string) error
}

type CreateProblemListParams struct {
	CreatorID   string   `json:"creator_id" binding:"required"`
	ProblemIDs  []string `json:"problem_ids" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Language    string   `json:"language" binding:"required,language"`
}

func listFromCreateParams(arg CreateProblemListParams) ProblemList {
	return ProblemList{
		CreatorID:   arg.CreatorID,
		ProblemIDs:  arg.ProblemIDs,
		Description: arg.Description,
		CreatedAt:   time.Now(),

		Upvotes:   0,
		Downvotes: 0,
	}
}

func (db *MongoDB) CreateProblemList(arg CreateProblemListParams) (ProblemList, error) {
	problemList := listFromCreateParams(arg)

	collection := db.client.Database("solvify").Collection("problemlists")
	result, err := collection.InsertOne(context.Background(), problemList)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problemList.ID = id

	return problemList, err
}

func (db *MongoDB) GetProblemList(id string) (ProblemList, error) {
	var problemList ProblemList
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return problemList, err
	}

	collection := db.client.Database("solvify").Collection("problemlists")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&problemList)
	return problemList, err
}

type ListProblemListsParams struct {
	Limit int32 `form:"limit"`
	Skip  int32 `form:"skip"`
}

func (db *MongoDB) ListProblemLists(arg ListProblemListsParams) ([]ProblemList, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(arg.Limit))
	findOptions.SetSkip(int64(arg.Skip))

	collection := db.client.Database("solvify").Collection("problemlists")
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	problemLists := make([]ProblemList, 0)
	for cursor.Next(context.Background()) {
		var list ProblemList
		err := cursor.Decode(&list)
		if err != nil {
			return nil, err
		}
		problemLists = append(problemLists, list)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return problemLists, err
}

type UpdateProblemListParams struct {
	ProblemIDs  []string `json:"problem_ids"`
	Description string   `json:"description"`
	Language    string   `json:"language" binding:"language"`
}

func (db *MongoDB) UpdateProblemList(arg UpdateProblemListParams, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("problemlists")
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"problem_ids": arg.ProblemIDs,
		"description": arg.Description,
		"language":    arg.Language,
	}}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	return err
}

func (db *MongoDB) DeleteProblemList(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("problemlists")
	filter := bson.D{{Key: "_id", Value: objectID}}
	result, err := collection.DeleteOne(context.Background(), filter)

	if result.DeletedCount == 0 {
		return errors.New("list not found")
	}

	return err
}
