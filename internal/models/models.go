package models

type Problem struct {
	ID        string `json:"_id" bson:"_id,omitempty"`
	Statement string `json:"statement"`
	Answer    string `json:"answer"`
}
