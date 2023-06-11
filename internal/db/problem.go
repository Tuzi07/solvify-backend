package db

import (
	"context"
	"errors"
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

	SolveTFProblem(arg SolveTFProblemParams) (TFProblemAttempt, error)
	SolveMTFProblem(arg SolveMTFProblemParams) (MTFProblemAttempt, error)
	SolveMCProblem(arg SolveMCProblemParams) (MCProblemAttempt, error)
	SolveMSProblem(arg SolveMSProblemParams) (MSProblemAttempt, error)

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
	Language         string `json:"language" binding:"required,language"`
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
	CorrectItems []bool   `json:"correct_items" binding:"required"`
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

func (db *MongoDB) CreateTFProblem(arg CreateTFProblemParams) (TFProblem, error) {
	problem := tfProblemFromCreateParams(arg)

	collection := db.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (db *MongoDB) CreateMTFProblem(arg CreateMTFProblemParams) (MTFProblem, error) {
	problem := mtfProblemFromCreateParams(arg)

	collection := db.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (db *MongoDB) CreateMCProblem(arg CreateMCProblemParams) (MCProblem, error) {
	problem := mcProblemFromCreateParams(arg)

	collection := db.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (db *MongoDB) CreateMSProblem(arg CreateMSProblemParams) (MSProblem, error) {
	problem := msProblemFromCreateParams(arg)

	collection := db.client.Database("solvify").Collection("problems")
	result, err := collection.InsertOne(context.Background(), problem)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	problem.ID = id

	return problem, err
}

func (db *MongoDB) GetProblem(id string) (AnyProblem, error) {
	var problem AnyProblem
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return problem, err
	}

	collection := db.client.Database("solvify").Collection("problems")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&problem)
	return problem, err
}

type UpdateProblemParams struct {
	Feedback         string `json:"feedback"`
	SubjectId        string `json:"subject_id"`
	TopicId          string `json:"topic_id"`
	SubtopicId       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education"`
	Language         string `json:"language" binding:"language"`
}

func (db *MongoDB) UpdateProblem(arg UpdateProblemParams, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("problems")
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

func (db *MongoDB) DeleteProblem(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("problems")
	filter := bson.D{{Key: "_id", Value: objectID}}
	result, err := collection.DeleteOne(context.Background(), filter)

	if result.DeletedCount == 0 {
		return errors.New("problem not found")
	}

	return err
}

type ListProblemsParams struct {
	Limit int32 `form:"limit"`
	Skip  int32 `form:"skip"`
}

func (db *MongoDB) ListProblems(arg ListProblemsParams) ([]AnyProblem, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(arg.Limit))
	findOptions.SetSkip(int64(arg.Skip))

	collection := db.client.Database("solvify").Collection("problems")
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

type SolveTFProblemParams struct {
	UserID       string `json:"user_id" binding:"required"`
	ProblemId    string `json:"problem_id" binding:"required"`
	BoolAnswer   *bool  `json:"bool_answer" binding:"required"`
	BoolResponse *bool  `json:"bool_response" binding:"required"`
}

func (db *MongoDB) SolveTFProblem(arg SolveTFProblemParams) (TFProblemAttempt, error) {
	attempt := tfProblemAttemptFromParams(arg)

	collection := db.client.Database("solvify").Collection("problem_attempts")
	result, err := collection.InsertOne(context.Background(), attempt)
	if err != nil {
		return attempt, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	attempt.ID = id

	collection = db.client.Database("solvify").Collection("user_problem_histories")
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemId}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":    arg.UserID,
			"problem_id": arg.ProblemId,
			"VoteStatus": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)

	return attempt, err
}

func tfProblemAttemptFromParams(arg SolveTFProblemParams) TFProblemAttempt {
	return TFProblemAttempt{
		ProblemAttempt: ProblemAttempt{
			AttemptedAt:      time.Now(),
			SolutionAccuracy: tfSolutionAccuracy(*arg.BoolAnswer, *arg.BoolResponse),
		},
		BoolResponse: *arg.BoolResponse,
	}
}

func tfSolutionAccuracy(answer bool, response bool) SolutionAccuracy {
	if answer == response {
		return Correct
	}
	return Incorrect
}

type SolveMTFProblemParams struct {
	UserID        string `json:"user_id" binding:"required"`
	ProblemId     string `json:"problem_id" binding:"required"`
	BoolAnswers   []bool `json:"bool_answers" binding:"required"`
	BoolResponses []bool `json:"bool_responses" binding:"required"`
}

func (db *MongoDB) SolveMTFProblem(arg SolveMTFProblemParams) (MTFProblemAttempt, error) {
	attempt := mtfProblemAttemptFromParams(arg)

	collection := db.client.Database("solvify").Collection("problem_attempts")
	result, err := collection.InsertOne(context.Background(), attempt)
	if err != nil {
		return attempt, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	attempt.ID = id

	collection = db.client.Database("solvify").Collection("user_problem_histories")
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemId}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":    arg.UserID,
			"problem_id": arg.ProblemId,
			"VoteStatus": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)

	return attempt, err
}

func mtfProblemAttemptFromParams(arg SolveMTFProblemParams) MTFProblemAttempt {
	return MTFProblemAttempt{
		ProblemAttempt: ProblemAttempt{
			AttemptedAt:      time.Now(),
			SolutionAccuracy: mtfSolutionAccuracy(arg.BoolAnswers, arg.BoolResponses),
		},
		BoolResponses: arg.BoolResponses,
	}
}

func mtfSolutionAccuracy(answers []bool, responses []bool) SolutionAccuracy {
	amountOfCorrectAnswers := 0
	amountOfItems := len(answers)

	for i := 0; i < amountOfItems; i++ {
		if answers[i] == responses[i] {
			amountOfCorrectAnswers++
		}
	}

	if amountOfCorrectAnswers == amountOfItems {
		return Correct
	} else if amountOfCorrectAnswers == 0 {
		return Incorrect
	} else {
		return Partial
	}
}

type SolveMCProblemParams struct {
	UserID       string `json:"user_id" binding:"required"`
	ProblemId    string `json:"problem_id" binding:"required"`
	CorrectItem  *int   `json:"correct_item" binding:"required"`
	ItemResponse *int   `json:"item_response" binding:"required"`
}

func (db *MongoDB) SolveMCProblem(arg SolveMCProblemParams) (MCProblemAttempt, error) {
	attempt := mcProblemAttemptFromParams(arg)

	collection := db.client.Database("solvify").Collection("problem_attempts")
	result, err := collection.InsertOne(context.Background(), attempt)
	if err != nil {
		return attempt, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	attempt.ID = id

	collection = db.client.Database("solvify").Collection("user_problem_histories")
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemId}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":    arg.UserID,
			"problem_id": arg.ProblemId,
			"VoteStatus": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)

	return attempt, err
}

func mcProblemAttemptFromParams(arg SolveMCProblemParams) MCProblemAttempt {
	return MCProblemAttempt{
		ProblemAttempt: ProblemAttempt{
			AttemptedAt:      time.Now(),
			SolutionAccuracy: mcSolutionAccuracy(*arg.CorrectItem, *arg.ItemResponse),
		},
		ItemResponse: *arg.ItemResponse,
	}
}

func mcSolutionAccuracy(answer int, response int) SolutionAccuracy {
	if answer == response {
		return Correct
	}
	return Incorrect
}

type SolveMSProblemParams struct {
	UserID        string `json:"user_id" binding:"required"`
	ProblemId     string `json:"problem_id" binding:"required"`
	CorrectItems  []bool `json:"correct_items" binding:"required"`
	ItemResponses []bool `json:"item_responses" binding:"required"`
}

func (db *MongoDB) SolveMSProblem(arg SolveMSProblemParams) (MSProblemAttempt, error) {
	attempt := msProblemAttemptFromParams(arg)

	collection := db.client.Database("solvify").Collection("problem_attempts")
	result, err := collection.InsertOne(context.Background(), attempt)
	if err != nil {
		return attempt, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	attempt.ID = id

	collection = db.client.Database("solvify").Collection("user_problem_histories")
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemId}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":    arg.UserID,
			"problem_id": arg.ProblemId,
			"VoteStatus": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)

	return attempt, err
}

func msProblemAttemptFromParams(arg SolveMSProblemParams) MSProblemAttempt {
	return MSProblemAttempt{
		ProblemAttempt: ProblemAttempt{
			AttemptedAt:      time.Now(),
			SolutionAccuracy: mtfSolutionAccuracy(arg.CorrectItems, arg.ItemResponses),
		},
		ItemResponses: arg.ItemResponses,
	}
}
