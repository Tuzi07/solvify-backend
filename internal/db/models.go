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
	SubjectID        string `json:"subject_id" bson:"subject_id"`
	TopicID          string `json:"topic_id" bson:"topic_id"`
	SubtopicID       string `json:"subtopic_id" bson:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education" bson:"level_of_education"`
	Language         string `json:"language" bson:"language"`

	Statement       string `json:"statement" bson:"statement"`
	CreatorID       string `json:"creator_id" bson:"creator_id"`
	CreatorUsername string `json:"creator_username" bson:"creator_username"`
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
	CorrectItems []bool   `json:"correct_items" bson:"correct_items"`
}

type AnyProblem struct {
	Problem      `bson:"inline"`
	BoolAnswer   bool     `json:"bool_answer" bson:"bool_answer,omitempty"`
	Items        []string `json:"items" bson:"items,omitempty"`
	BoolAnswers  []bool   `json:"bool_answers" bson:"bool_answers,omitempty"`
	CorrectItem  int      `json:"correct_item" bson:"correct_item,omitempty"`
	CorrectItems []int    `json:"correct_items" bson:"correct_items,omitempty"`
}

type User struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	Email       string    `json:"email" bson:"email"`
	Username    string    `json:"username" bson:"username"`
	DisplayName string    `json:"display_name" bson:"display_name"`
	Languages   []string  `json:"languages" bson:"languages"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`

	CreatedProblems    []string `json:"created_problems" bson:"created_problems"`
	CreatedLists       []string `json:"created_lists" bson:"created_lists"`
	SolveLaterProblems []string `json:"solve_later_problems" bson:"solve_later_problems"`
	SolveLaterLists    []string `json:"solve_later_lists" bson:"solve_later_lists"`
}

type VoteStatus int

const (
	Downvote VoteStatus = iota
	NoVote
	Upvote
)

type UserProblemHistory struct {
	ID                 string     `json:"_id" bson:"_id,omitempty"`
	UserID             string     `json:"user_id" bson:"user_id"`
	ProblemID          string     `json:"problem_id" bson:"problem_id"`
	ProblemAttemptsIDs []string   `json:"problem_attempts_ids" bson:"problem_attempts_ids"`
	VoteStatus         VoteStatus `json:"vote_status" bson:"vote_status"`
}

type SolutionAccuracy int

const (
	Incorrect SolutionAccuracy = iota
	Partial
	Correct
)

type ProblemAttempt struct {
	ID               string           `json:"_id" bson:"_id,omitempty"`
	AttemptedAt      time.Time        `json:"attempted_at" bson:"attempted_at"`
	SolutionAccuracy SolutionAccuracy `json:"solution_accuracy" bson:"solution_accuracy"`
}

type TFProblemAttempt struct {
	ProblemAttempt `bson:"inline"`
	BoolResponse   bool `json:"bool_response" bson:"bool_response"`
}

type MTFProblemAttempt struct {
	ProblemAttempt `bson:"inline"`
	BoolResponses  []bool `json:"bool_responses" bson:"bool_responses"`
}

type MCProblemAttempt struct {
	ProblemAttempt `bson:"inline"`
	ItemResponse   int `json:"item_response" bson:"item_response"`
}

type MSProblemAttempt struct {
	ProblemAttempt `bson:"inline"`
	ItemResponses  []bool `json:"item_responses" bson:"item_responses"`
}

type ProblemList struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	CreatorID   string    `json:"creator_id" bson:"creator_id"`
	ProblemIDs  []string  `json:"problem_ids" bson:"problem_ids"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`

	Upvotes   int `json:"upvotes" bson:"upvotes"`
	Downvotes int `json:"downvotes" bson:"downvotes"`
}

type Subject struct {
	ID       string `json:"_id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Language string `json:"language" bson:"language"`
}

type Topic struct {
	ID        string `json:"_id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name"`
	SubjectID string `json:"subject_id" bson:"subject_id"`
}

type Subtopic struct {
	ID      string `json:"_id" bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name"`
	TopicID string `json:"topic_id" bson:"topic_id"`
}

type ProblemReport struct {
	ID         string    `json:"_id" bson:"_id,omitempty"`
	ProblemID  string    `json:"problem_id" bson:"problem_id"`
	UserID     string    `json:"user_id" bson:"user_id"`
	ReportedAt time.Time `json:"reported_at" bson:"reported_at"`
	Reason     string    `json:"reason" bson:"reason"`
}

type ListReport struct {
	ID         string    `json:"_id" bson:"_id,omitempty"`
	ListID     string    `json:"list_id" bson:"list_id"`
	UserID     string    `json:"user_id" bson:"user_id"`
	ReportedAt time.Time `json:"reported_at" bson:"reported_at"`
	Reason     string    `json:"reason" bson:"reason"`
}

type ProblemEditSuggestion struct {
	ID              string    `json:"_id" bson:"_id,omitempty"`
	ProblemID       string    `json:"problem_id" bson:"problem_id"`
	CreatorUsername string    `json:"user_id" bson:"user_id"`
	SuggestedAt     time.Time `json:"suggested_at" bson:"suggested_at"`

	Feedback         string `json:"feedback" bson:"feedback"`
	SubjectID        string `json:"subject_id" bson:"subject_id"`
	TopicID          string `json:"topic_id" bson:"topic_id"`
	SubtopicID       string `json:"subtopic_id" bson:"subtopic_id"`
	LevelOfEducation string `json:"level_of_education" bson:"level_of_education"`
	Language         string `json:"language" bson:"language"`
}
