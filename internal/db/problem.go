package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProblemDatabase interface {
	CreateTFProblem(arg CreateTFProblemParams) (TFProblem, error)
	CreateMTFProblem(arg CreateMTFProblemParams) (MTFProblem, error)
	CreateMCProblem(arg CreateMCProblemParams) (MCProblem, error)
	CreateMSProblem(arg CreateMSProblemParams) (MSProblem, error)

	GetProblem(id string) (AnyProblem, error)
	UpdateProblem(arg UpdateProblemParams, id string) error
	DeleteProblem(id string) error
	ListProblems(arg ListProblemsParams) ([]AnyProblem, error)
}

type CreateProblemParams struct {
	Statement        string `json:"statement" binding:"required"`
	Feedback         string `json:"feedback"`
	SubjectId        string `json:"subject_id"`
	TopicId          string `json:"topic_id"`
	SubtopicId       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education"`
	Language         string `json:"language" binding:"required,oneof=en pt"`
	CreatorID        string `json:"creator_id" binding:"required"`
}

type CreateTFProblemParams struct {
	CreateProblemParams
	BoolAnswer *bool `json:"bool_answer" binding:"required"` // pointer explained here: https://github.com/go-playground/validator/issues/692
}

type CreateMTFProblemParams struct {
	CreateProblemParams
	Items       []string `json:"items" binding:"required"`
	BoolAnswers []bool   `json:"bool_answers" binding:"required"`
}

type CreateMCProblemParams struct {
	CreateProblemParams
	Items       []string `json:"items" binding:"required"`
	CorrectItem *int     `json:"correct_item" binding:"required"`
}

type CreateMSProblemParams struct {
	CreateProblemParams
	Items        []string `json:"items" binding:"required"`
	CorrectItems []int    `json:"correct_items" binding:"required"`
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
		BoolAnswer: *arg.BoolAnswer,
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
		CorrectItem: *arg.CorrectItem,
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
	Limit int32 `form:"limit"`
	Skip  int32 `form:"skip"`
}

func (dbManager *MongoDB) ListProblems(arg ListProblemsParams) ([]AnyProblem, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(arg.Limit))
	findOptions.SetSkip(int64(arg.Skip))

	collection := dbManager.client.Database("solvify").Collection("problems")
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	problems := make([]AnyProblem, 0)
	for cursor.Next(context.Background()) {
		var problem AnyProblem
		err := cursor.Decode(&problem)
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return problems, err
}
