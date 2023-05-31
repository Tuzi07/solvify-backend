package db

import "time"

type ProblemType int

const (
	TrueFalse ProblemType = iota
	MultipleTrueFalse
	MultipleChoice
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

// TFProblem represents a True False Problem
type TFProblem struct {
	Problem    `bson:"inline"`
	BoolAnswer bool `json:"bool_answer" bson:"bool_answer"`
}

// MTFProblem represents a Multiple True False Problem
type MTFProblem struct {
	Problem     `bson:"inline"`
	Items       []string `json:"items" bson:"items"`
	BoolAnswers []bool   `json:"bool_answers" bson:"bool_answers"`
}

// MCProblem represents a Multiple Choice Problem
type MCProblem struct {
	Problem     `bson:"inline"`
	Items       []string `json:"items" bson:"items"`
	CorrectItem int      `json:"correct_item" bson:"correct_item"`
}

// MSProblem represents a Multiple Selection Problem
type MSProblem struct {
	Problem      `bson:"inline"`
	Items        []string `json:"items" bson:"items"`
	CorrectItems []int    `json:"correct_items" bson:"correct_items"`
}

type ProblemList struct {
	ID string `json:"_id" bson:"_id,omitempty"`

	ProblemIds []string `json:"problem_ids" bson:"problem_ids"`

	Description string `json:"description" bson:"description"`

	CreatorID string    `json:"creator_id" bson:"creator_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	Upvotes   int `json:"upvotes" bson:"upvotes"`
	Downvotes int `json:"downvotes" bson:"downvotes"`
}

type User struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	Username    string    `json:"username" bson:"username"`
	DisplayName string    `json:"display_name" bson:"display_name"`
	Languages   []string  `json:"languages" bson:"languages"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`

	CreatedProblems    []string `json:"created_problems" bson:"created_problems"`
	CreatedLists       []string `json:"created_lists" bson:"created_lists"`
	SolveLaterProblems []string `json:"solve_later_problems" bson:"solve_later_problems"`
	SolveLaterLists    []string `json:"solve_later_lists" bson:"solve_later_lists"`
}
