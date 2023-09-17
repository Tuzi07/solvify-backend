package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProblemAttemptDatabase interface {
	ListUserAttempts(id string, pagination PaginationParams) ([]ProblemAttemptTableRow, error)
}

type ProblemAttemptTableRow struct {
	Subject          string           `json:"subject"`
	Statement        string           `json:"statement"`
	ProblemID        string           `json:"problem_id"`
	AttemptedAt      time.Time        `json:"attempted_at"`
	SolutionAccuracy SolutionAccuracy `json:"solution_accuracy"`
}

func (db *MongoDB) ListUserAttempts(id string, pagination PaginationParams) ([]ProblemAttemptTableRow, error) {
	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("problem_attempts")
	filter := bson.M{"user_id": id}
	findOptions := options.Find()
	findOptions.SetLimit(int64(pagination.Limit))
	findOptions.SetSkip(int64(pagination.Skip))

	descending := -1
	findOptions.SetSort(bson.D{{Key: "attempted_at", Value: descending}})

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	attempts := make([]ProblemAttemptTableRow, 0)
	for cursor.Next(context.Background()) {
		var problemAttempt ProblemAttempt
		if err := cursor.Decode(&problemAttempt); err != nil {
			return nil, err
		}

		attempt, err := db.attemptTableRowFromAttempt(problemAttempt)
		if err != nil {
			return nil, err
		}

		attempts = append(attempts, attempt)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return attempts, nil
}

func (db *MongoDB) attemptTableRowFromAttempt(attempt ProblemAttempt) (ProblemAttemptTableRow, error) {
	var attemptTableRow ProblemAttemptTableRow

	id := attempt.ProblemID
	problem, err := db.GetProblem(id)
	if err != nil {
		return attemptTableRow, err
	}

	subject, err := db.GetSubject(problem.SubjectID)
	if err != nil {
		return attemptTableRow, err
	}

	return ProblemAttemptTableRow{
		Subject:          subject.Name,
		Statement:        problem.Statement,
		ProblemID:        id,
		AttemptedAt:      attempt.AttemptedAt,
		SolutionAccuracy: attempt.SolutionAccuracy,
	}, nil
}
