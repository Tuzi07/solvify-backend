package models

type Problem struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
	Answer    string `json:"answer"`
}
