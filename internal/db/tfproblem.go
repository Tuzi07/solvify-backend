package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TFProblemDatabase interface {
	CreateTFProblem(arg CreateTFProblemParams) (TFProblem, error)
	GetTFProblem(id string) (TFProblem, error)
	UpdateTFProblem(arg UpdateTFProblemParams, id string) error
	DeleteTFProblem(id string) error
}

type CreateTFProblemParams struct {
	Statement        string `json:"statement"`
	Feedback         string `json:"feedback"`
	SubjectId        string `json:"subject_id"`
	TopicId          string `json:"topic_id"`
	SubtopicId       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education"`
	Language         string `json:"language"`
	CreatorID        string `json:"creator_id"`
	BoolAnswer       bool   `json:"bool_answer"`
}

func problemFromCreateParams(arg CreateTFProblemParams) TFProblem {
	return TFProblem{
		Problem: Problem{
			Statement:        arg.Statement,
			Feedback:         arg.Feedback,
			SubjectId:        arg.SubjectId,
			TopicId:          arg.TopicId,
			SubtopicId:       arg.SubtopicId,
			LevelOfEducation: arg.LevelOfEducation,
			Language:         arg.Language,
			CreatorID:        arg.CreatorID,

			ProblemType:    TrueFalse,
			CreatedAt:      time.Now(),
			Attempts:       0,
			CorrectAnswers: 0,
			Accuracy:       0.0,
			Upvotes:        0,
			Downvotes:      0,
		},
		BoolAnswer: arg.BoolAnswer,
	}
}

func (dbManager *MongoDB) CreateTFProblem(arg CreateTFProblemParams) (TFProblem, error) {
	problem := problemFromCreateParams(arg)

	collection := dbManager.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (dbManager *MongoDB) GetTFProblem(id string) (TFProblem, error) {
	var problem TFProblem
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return problem, err
	}

	collection := dbManager.client.Database("solvify").Collection("problems")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&problem)
	if err != nil {
		return problem, err
	}
	return problem, nil
}

type UpdateTFProblemParams struct {
	Feedback         string `json:"feedback"`
	SubjectId        string `json:"subject_id"`
	TopicId          string `json:"topic_id"`
	SubtopicId       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education"`
	Language         string `json:"language"`
}

func (dbManager *MongoDB) UpdateTFProblem(arg UpdateTFProblemParams, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := dbManager.client.Database("solvify").Collection("problems")
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"feedback":           arg.Feedback,
		"subject_id":         arg.SubjectId,
		"topic_id":           arg.TopicId,
		"subtopic_id":        arg.SubtopicId,
		"level_of_education": arg.LevelOfEducation,
		"language":           arg.Language,
	}}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	return err
}

func (dbManager *MongoDB) DeleteTFProblem(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := dbManager.client.Database("solvify").Collection("problems")
	filter := bson.D{{Key: "_id", Value: objectID}}
	_, err = collection.DeleteOne(context.Background(), filter)

	return err
}
