package db

import "time"

type ProblemType int

const (
	TrueFalse ProblemType = iota
	MultipleTrueFalse
	Multichoice
	MultipleSelection
)

type Problem struct {
	ID             string      `json:"_id" bson:"_id,omitempty"`
	ProblemType    ProblemType `json:"problem_type" bson:"problem_type"`
	CreatedAt      time.Time   `json:"created_at" bson:"created_at"`
	Attempts       int         `json:"attempts" bson:"attempts"`
	CorrectAnswers int         `json:"correct_answers" bson:"correct_answers"`
	Accuracy       float64     `json:"accuracy" bson:"accuracy"`
	Upvotes        int         `json:"upvotes" bson:"upvotes"`
	Downvotes      int         `json:"downvotes" bson:"downvotes"`

	Feedback         string `json:"feedback" bson:"feedback"`
	SubjectId        string `json:"subject_id" bson:"subject_id"`
	TopicId          string `json:"topic_id" bson:"topic_id"`
	SubtopicId       string `json:"subtopic_id" bson:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education" bson:"level_of_education"`
	Language         string `json:"language" bson:"language"`

	Statement string `json:"statement" bson:"statement"`
	CreatorID string `json:"creator_id" bson:"creator_id"`
}

// True False Problem
type TFProblem struct {
	Problem    `bson:"inline"`
	BoolAnswer bool `json:"bool_answer" bson:"bool_answer"`
}

// Multiple True False Problem
type MTFProblem struct {
	Problem     `bson:"inline"`
	Items       []string `json:"items" bson:"items"`
	BoolAnswers []bool   `json:"bool_answers" bson:"bool_answers"`
}

// Multiple Choice Problem
type MCProblem struct {
	Problem     `bson:"inline"`
	Items       []string `json:"items" bson:"items"`
	CorrectItem int      `json:"correct_item" bson:"correct_item"`
}

// Multiple Selection Problem
type MSProblem struct {
	Problem      `bson:"inline"`
	Items        []string `json:"items" bson:"items"`
	CorrectItems []int    `json:"correct_items" bson:"correct_items"`
}

type ProblemList struct {
	ID string `json:"_id" bson:"_id"`

	ProblemIds []string `json:"problem_ids" bson:"problem_ids"`

	Description string `json:"description" bson:"description"`

	CreatorID string    `json:"creator_id" bson:"creator_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	Upvotes   int `json:"upvotes" bson:"upvotes"`
	Downvotes int `json:"downvotes" bson:"downvotes"`
}
