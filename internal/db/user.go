package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDatabase interface {
	CreateUser(arg CreateUserParams) (User, error)
	GetUser(id string) (User, error)
	UpdateUser(arg UpdateUserParams, id string) error
	DeleteUser(id string) error
	CheckEmailUnique(arg CheckEmailUniqueParams) (bool, error)
}

type CreateUserParams struct {
	Username    string   `json:"username" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	DisplayName string   `json:"display_name" binding:"required"`
	Languages   []string `json:"languages" binding:"required,languages"`
}

func UserFromParams(arg CreateUserParams) User {
	return User{
		Username:    arg.Username,
		Email:       arg.Email,
		DisplayName: arg.DisplayName,
		Languages:   arg.Languages,
		CreatedAt:   time.Now(),

		CreatedProblems:    []string{},
		CreatedLists:       []string{},
		SolveLaterProblems: []string{},
		SolveLaterLists:    []string{},
	}
}

func (db *MongoDB) CreateUser(arg CreateUserParams) (User, error) {
	user := UserFromParams(arg)

	if !db.isUniqueUsername(user.Username) {
		return user, errors.New("username already exists")
	}

	collection := db.client.Database("solvify").Collection("users")
	result, err := collection.InsertOne(context.Background(), user)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	user.ID = id

	return user, err
}

func (db *MongoDB) isUniqueUsername(username string) bool {
	collection := db.client.Database("solvify").Collection("users")
	filter := bson.D{{Key: "username", Value: username}}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	return count == 0
}

func (db *MongoDB) GetUser(id string) (User, error) {
	var user User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	collection := db.client.Database("solvify").Collection("users")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&user)

	return user, err
}

type UpdateUserParams struct {
	DisplayName string   `json:"display_name"`
	Languages   []string `json:"languages" binding:"languages"`
}

func (db *MongoDB) UpdateUser(arg UpdateUserParams, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("users")
	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.M{"$set": bson.M{
		"display_name": arg.DisplayName,
		"languages":    arg.Languages,
	}}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return err
}

func (db *MongoDB) DeleteUser(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.client.Database("solvify").Collection("users")
	filter := bson.D{{Key: "_id", Value: objectID}}
	result, err := collection.DeleteOne(context.Background(), filter)

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return err
}

type CheckEmailUniqueParams struct {
	Email string `json:"email" binding:"required,email"`
}

func (db *MongoDB) CheckEmailUnique(arg CheckEmailUniqueParams) (bool, error) {
	collection := db.client.Database("solvify").Collection("users")
	filter := bson.D{{Key: "email", Value: arg.Email}}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count == 0, err
}
