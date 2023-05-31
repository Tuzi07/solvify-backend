package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProblemDatabase interface {
	CreateTFProblem(arg CreateTFProblemParams) (TFProblem, error)
	CreateMTFProblem(arg CreateMTFProblemParams) (MTFProblem, error)
	CreateMCProblem(arg CreateMCProblemParams) (MCProblem, error)
	CreateMSProblem(arg CreateMSProblemParams) (MSProblem, error)

	GetProblem(id string) (AnyProblem, error)
	UpdateProblem(arg UpdateProblemParams, id string) error
	DeleteProblem(id string) error
	ListProblems(arg ListProblemsParams) ([]Problem, error)
}

type CreateProblemParams struct {
	Statement        string `json:"statement"`
	Feedback         string `json:"feedback"`
	SubjectId        string `json:"subject_id"`
	TopicId          string `json:"topic_id"`
	SubtopicId       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education"`
	Language         string `json:"language"`
	CreatorID        string `json:"creator_id"`
}

type CreateTFProblemParams struct {
	CreateProblemParams
	BoolAnswer bool `json:"bool_answer"`
}

type CreateMTFProblemParams struct {
	CreateProblemParams
	Items       []string `json:"items"`
	BoolAnswers []bool   `json:"bool_answers"`
}

type CreateMCProblemParams struct {
	CreateProblemParams
	Items       []string `json:"items"`
	CorrectItem int      `json:"correct_item"`
}

type CreateMSProblemParams struct {
	CreateProblemParams
	Items        []string `json:"items"`
	CorrectItems []int    `json:"correct_items"`
}

func problemFromCreateParams(arg CreateProblemParams, problemType ProblemType) Problem {
	return Problem{
		Statement:        arg.Statement,
		Feedback:         arg.Feedback,
		SubjectId:        arg.SubjectId,
		TopicId:          arg.TopicId,
		SubtopicId:       arg.SubtopicId,
		LevelOfEducation: arg.LevelOfEducation,
		Language:         arg.Language,
		CreatorID:        arg.CreatorID,

		ProblemType:    problemType,
		CreatedAt:      time.Now(),
		Attempts:       0,
		CorrectAnswers: 0,
		Accuracy:       0.0,
		Upvotes:        0,
		Downvotes:      0,
	}
}

func tfProblemFromCreateParams(arg CreateTFProblemParams) TFProblem {
	return TFProblem{
		Problem:    problemFromCreateParams(arg.CreateProblemParams, TrueFalse),
		BoolAnswer: arg.BoolAnswer,
	}
}

func mtfProblemFromCreateParams(arg CreateMTFProblemParams) MTFProblem {
	return MTFProblem{
		Problem:     problemFromCreateParams(arg.CreateProblemParams, MultipleTrueFalse),
		Items:       arg.Items,
		BoolAnswers: arg.BoolAnswers,
	}
}

func mcProblemFromCreateParams(arg CreateMCProblemParams) MCProblem {
	return MCProblem{
		Problem:     problemFromCreateParams(arg.CreateProblemParams, MultipleChoice),
		Items:       arg.Items,
		CorrectItem: arg.CorrectItem,
	}
}

func msProblemFromCreateParams(arg CreateMSProblemParams) MSProblem {
	return MSProblem{
		Problem:      problemFromCreateParams(arg.CreateProblemParams, MultipleSelection),
		Items:        arg.Items,
		CorrectItems: arg.CorrectItems,
	}
}

func (dbManager *MongoDB) CreateTFProblem(arg CreateTFProblemParams) (TFProblem, error) {
	problem := tfProblemFromCreateParams(arg)

	collection := dbManager.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (dbManager *MongoDB) CreateMTFProblem(arg CreateMTFProblemParams) (MTFProblem, error) {
	problem := mtfProblemFromCreateParams(arg)

	collection := dbManager.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (dbManager *MongoDB) CreateMCProblem(arg CreateMCProblemParams) (MCProblem, error) {
	problem := mcProblemFromCreateParams(arg)

	collection := dbManager.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (dbManager *MongoDB) CreateMSProblem(arg CreateMSProblemParams) (MSProblem, error) {
	problem := msProblemFromCreateParams(arg)

	collection := dbManager.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (dbManager *MongoDB) GetProblem(id string) (AnyProblem, error) {
	var problem AnyProblem
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

type UpdateProblemParams struct {
	Feedback         string `json:"feedback"`
	SubjectId        string `json:"subject_id"`
	TopicId          string `json:"topic_id"`
	SubtopicId       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education"`
	Language         string `json:"language"`
}

func (dbManager *MongoDB) UpdateProblem(arg UpdateProblemParams, id string) error {
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

func (dbManager *MongoDB) DeleteProblem(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := dbManager.client.Database("solvify").Collection("problems")
	filter := bson.D{{Key: "_id", Value: objectID}}
	_, err = collection.DeleteOne(context.Background(), filter)

	return err
}

type ListProblemsParams struct {
}

func (dbManager *MongoDB) ListProblems(arg ListProblemsParams) ([]Problem, error) {
	return nil, nil
}
