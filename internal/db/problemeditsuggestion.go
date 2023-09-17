package db

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProblemEditSuggestionDatabase interface {
	CreateProblemEditSuggestion(arg CreateProblemEditSuggestionParams) (ProblemEditSuggestion, error)
	GetProblemEditSuggestion(id string) (ProblemEditSuggestion, error)
	DeleteProblemEditSuggestion(id string) error

	ListSuggestionsOfProblem(problemID string) ([]ProblemEditSuggestion, error)
	AcceptProblemEditSuggestion(id string) (AnyProblem, error)
}

type CreateProblemEditSuggestionParams struct {
	ProblemID       string `json:"problem_id" binding:"required"`
	CreatorUsername string `json:"creator_username" binding:"required"`
	UpdateProblemParams
}

func suggestionFromParams(arg CreateProblemEditSuggestionParams) ProblemEditSuggestion {
	return ProblemEditSuggestion{
		ProblemID:       arg.ProblemID,
		CreatorUsername: arg.CreatorUsername,
		SuggestedAt:     time.Now(),

		Feedback:         arg.Feedback,
		SubjectID:        arg.SubjectID,
		TopicID:          arg.TopicID,
		SubtopicID:       arg.SubtopicID,
		LevelOfEducation: arg.LevelOfEducation,
		Language:         arg.Language,
	}
}

func (db *MongoDB) CreateProblemEditSuggestion(arg CreateProblemEditSuggestionParams) (ProblemEditSuggestion, error) {
	suggestion := suggestionFromParams(arg)

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("problem_edit_suggestions")
	result, err := collection.InsertOne(context.Background(), suggestion)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	suggestion.ID = id

	return suggestion, err
}

func (db *MongoDB) GetProblemEditSuggestion(id string) (ProblemEditSuggestion, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ProblemEditSuggestion{}, err
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("problem_edit_suggestions")
	filter := bson.M{"_id": objectID}
	result := collection.FindOne(context.Background(), filter)

	var suggestion ProblemEditSuggestion
	err = result.Decode(&suggestion)
	return suggestion, err
}

func (db *MongoDB) DeleteProblemEditSuggestion(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("problem_edit_suggestions")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})

	if result.DeletedCount == 0 {
		return errors.New("problem edit suggestion not found")
	}

	return err
}

func (db *MongoDB) ListSuggestionsOfProblem(problemID string) ([]ProblemEditSuggestion, error) {
	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("problem_edit_suggestions")
	filter := bson.M{"problem_id": problemID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var suggestions []ProblemEditSuggestion
	for cursor.Next(context.Background()) {
		var suggestion ProblemEditSuggestion
		err := cursor.Decode(&suggestion)
		if err != nil {
			return nil, err
		}
		suggestions = append(suggestions, suggestion)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return suggestions, nil
}

func (db *MongoDB) AcceptProblemEditSuggestion(id string) (AnyProblem, error) {
	var problem AnyProblem
	suggestion, err := db.GetProblemEditSuggestion(id)
	if err != nil {
		return problem, err
	}

	problem, err = db.GetProblem(suggestion.ProblemID)
	if err != nil {
		return problem, err
	}

	problem = editedProblem(problem, suggestion)

	params := UpdateProblemParams{
		Feedback:         problem.Feedback,
		SubjectID:        problem.SubjectID,
		TopicID:          problem.TopicID,
		SubtopicID:       problem.SubtopicID,
		LevelOfEducation: problem.LevelOfEducation,
		Language:         problem.Language,
	}

	err = db.UpdateProblem(params, suggestion.ProblemID)
	if err != nil {
		return problem, err
	}

	err = db.DeleteProblemEditSuggestion(id)
	if err != nil {
		return problem, err
	}

	return problem, nil
}

func editedProblem(problem AnyProblem, suggestion ProblemEditSuggestion) AnyProblem {
	if suggestion.Feedback != "" {
		problem.Feedback = suggestion.Feedback
	}
	if suggestion.SubjectID != "" {
		problem.SubjectID = suggestion.SubjectID
	}
	if suggestion.TopicID != "" {
		problem.TopicID = suggestion.TopicID
	}
	if suggestion.SubtopicID != "" {
		problem.SubtopicID = suggestion.SubtopicID
	}
	if suggestion.LevelOfEducation != "" {
		problem.LevelOfEducation = suggestion.LevelOfEducation
	}
	if suggestion.Language != "" {
		problem.Language = suggestion.Language
	}

	return problem
}
