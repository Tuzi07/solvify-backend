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

	VoteProblem(arg VoteProblemParams) error
	CreateProblemReport(arg ReportProblemParams) (ProblemReport, error)
}

type CreateProblemParams struct {
	Statement        string `json:"statement" binding:"required"`
	Feedback         string `json:"feedback"`
	SubjectID        string `json:"subject_id"`
	TopicID          string `json:"topic_id"`
	SubtopicID       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education" binding:"level_of_education"`
	Language         string `json:"language" binding:"required,language"`
	CreatorID        string `json:"creator_id" binding:"required"`
	CreatorUsername  string `json:"creator_username" binding:"required"`
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
		SubjectID:        arg.SubjectID,
		TopicID:          arg.TopicID,
		SubtopicID:       arg.SubtopicID,
		LevelOfEducation: arg.LevelOfEducation,
		Language:         arg.Language,
		CreatorID:        arg.CreatorID,
		CreatorUsername:  arg.CreatorUsername,

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
	SubjectID        string `json:"subject_id"`
	TopicID          string `json:"topic_id"`
	SubtopicID       string `json:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education" binding:"level_of_education"`
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
		"subject_id":         arg.SubjectID,
		"topic_id":           arg.TopicID,
		"subtopic_id":        arg.SubtopicID,
		"level_of_education": arg.LevelOfEducation,
		"language":           arg.Language,
	}}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	if result.MatchedCount == 0 {
		return errors.New("problem not found")
	}

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

type PaginationParams struct {
	Limit int32 `form:"limit"`
	Skip  int32 `form:"skip"`
}

type ListProblemsParams struct {
	PaginationParams

	OrderBy    string `json:"order_by" binding:"field_to_order_problems"`
	Descending *bool  `json:"descending"`

	ProblemTypeFilter      *int   `json:"problem_type"`
	SubjectFilter          string `json:"subject_id"`
	TopicFilter            string `json:"topic_id"`
	SubtopicFilter         string `json:"subtopic_id"`
	LevelOfEducationFilter string `json:"level_of_education" binding:"level_of_education"`
	LanguageFilter         string `json:"language" binding:"language"`
}

// ListProblems returns a list of problems.
// The returned list is ordered by the field specified in the `order_by` parameter.
// order_by can be one of the following values: "created_at", "attempts", "accuracy", "upvotes"
// The Filter can be empty. If a filter is empty, it is ignored.
func (db *MongoDB) ListProblems(arg ListProblemsParams) ([]AnyProblem, error) {
	collection := db.client.Database("solvify").Collection("problems")
	filter := filterFromParams(arg)
	findOptions := optionsFromParams(arg)

	cursor, err := collection.Find(context.Background(), filter, findOptions)
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

func filterFromParams(arg ListProblemsParams) bson.M {
	filter := bson.M{}

	if arg.ProblemTypeFilter != nil {
		filter["problem_type"] = arg.ProblemTypeFilter
	}
	if arg.SubjectFilter != "" {
		filter["subject_id"] = arg.SubjectFilter
	}
	if arg.TopicFilter != "" {
		filter["topic_id"] = arg.TopicFilter
	}
	if arg.SubtopicFilter != "" {
		filter["subtopic_id"] = arg.SubtopicFilter
	}
	if arg.LevelOfEducationFilter != "" {
		filter["level_of_education"] = arg.LevelOfEducationFilter
	}
	if arg.LanguageFilter != "" {
		filter["language"] = arg.LanguageFilter
	}

	return filter
}

func optionsFromParams(arg ListProblemsParams) *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetLimit(int64(arg.Limit))
	findOptions.SetSkip(int64(arg.Skip))

	if arg.OrderBy == "" {
		arg.OrderBy = "created_at"
	}
	if arg.Descending == nil {
		descending := true
		arg.Descending = &descending
	}

	var sort bson.D
	if *arg.Descending {
		sort = bson.D{{Key: arg.OrderBy, Value: -1}}
	} else {
		sort = bson.D{{Key: arg.OrderBy, Value: 1}}
	}
	findOptions.SetSort(sort)

	return findOptions
}

type SolveProblemParams struct {
	UserID    string `json:"user_id" binding:"required"`
	ProblemID string `json:"problem_id" binding:"required"`
}

type SolveTFProblemParams struct {
	SolveProblemParams
	BoolAnswer   *bool `json:"bool_answer" binding:"required"`
	BoolResponse *bool `json:"bool_response" binding:"required"`
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
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemID}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":     arg.UserID,
			"problem_id":  arg.ProblemID,
			"vote_status": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return attempt, err
	}

	err = db.updateProblemAttempts(arg.SolveProblemParams, attempt.SolutionAccuracy)

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

func (db *MongoDB) updateProblemAttempts(arg SolveProblemParams, solutionAccuracy SolutionAccuracy) error {
	problem, err := db.GetProblem(arg.ProblemID)
	if err != nil {
		return err
	}

	attempts := problem.Attempts + 1
	correctAnswers := problem.CorrectAnswers
	if solutionAccuracy == Correct {
		correctAnswers++
	}
	accuracy := float32(correctAnswers) / float32(attempts)

	objectID, err := primitive.ObjectIDFromHex(arg.ProblemID)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("problems")
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"attempts":        attempts,
		"correct_answers": correctAnswers,
		"accuracy":        accuracy,
	}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	return err
}

type SolveMTFProblemParams struct {
	SolveProblemParams
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
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemID}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":     arg.UserID,
			"problem_id":  arg.ProblemID,
			"vote_status": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return attempt, err
	}

	err = db.updateProblemAttempts(arg.SolveProblemParams, attempt.SolutionAccuracy)

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
	SolveProblemParams
	CorrectItem  *int `json:"correct_item" binding:"required"`
	ItemResponse *int `json:"item_response" binding:"required"`
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
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemID}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":     arg.UserID,
			"problem_id":  arg.ProblemID,
			"vote_status": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return attempt, err
	}

	err = db.updateProblemAttempts(arg.SolveProblemParams, attempt.SolutionAccuracy)

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
	SolveProblemParams
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
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemID}
	update := bson.M{
		"$push": bson.M{"problem_attempts_ids": id},
		"$setOnInsert": bson.M{
			"user_id":     arg.UserID,
			"problem_id":  arg.ProblemID,
			"vote_status": NoVote,
		},
	}
	options := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return attempt, err
	}

	err = db.updateProblemAttempts(arg.SolveProblemParams, attempt.SolutionAccuracy)

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

type VoteProblemParams struct {
	UserID     string      `json:"user_id" binding:"required"`
	ProblemID  string      `json:"problem_id" binding:"required"`
	VoteStatus *VoteStatus `json:"vote_status" binding:"required"`
}

func (db *MongoDB) VoteProblem(arg VoteProblemParams) error {
	var userProblemHistory UserProblemHistory
	collection := db.client.Database("solvify").Collection("user_problem_histories")
	filter := bson.M{"user_id": arg.UserID, "problem_id": arg.ProblemID}
	err := collection.FindOne(context.Background(), filter).Decode(&userProblemHistory)
	if err != nil {
		return err
	}

	if userProblemHistory.VoteStatus == *arg.VoteStatus {
		return nil
	}

	update := bson.M{"$set": bson.M{"vote_status": *arg.VoteStatus}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	collection = db.client.Database("solvify").Collection("problems")
	objectID, err := primitive.ObjectIDFromHex(arg.ProblemID)
	if err != nil {
		return err
	}

	filter = bson.M{"_id": objectID}
	if userProblemHistory.VoteStatus == Upvote {
		update = bson.M{"$inc": bson.M{"upvotes": -1}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
	} else if userProblemHistory.VoteStatus == Downvote {
		update = bson.M{"$inc": bson.M{"downvotes": -1}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
	}

	if *arg.VoteStatus == Upvote {
		update = bson.M{"$inc": bson.M{"upvotes": 1}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
	} else if *arg.VoteStatus == Downvote {
		update = bson.M{"$inc": bson.M{"downvotes": 1}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
	}

	return err
}

type ReportProblemParams struct {
	ProblemID string `json:"problem_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

func reportFromParams(arg ReportProblemParams) ProblemReport {
	return ProblemReport{
		ProblemID:  arg.ProblemID,
		UserID:     arg.UserID,
		Reason:     arg.Reason,
		ReportedAt: time.Now(),
	}
}

func (db *MongoDB) CreateProblemReport(arg ReportProblemParams) (ProblemReport, error) {
	report := reportFromParams(arg)

	collection := db.client.Database("solvify").Collection("problem_reports")
	result, err := collection.InsertOne(context.Background(), report)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	report.ID = id

	return report, err
}
