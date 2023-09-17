package db

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LabelsDatabase interface {
	CreateSubject(arg CreateSubjectParams) (Subject, error)
	GetSubject(id string) (Subject, error)
	ListSubjects() ([]Subject, error)

	CreateTopic(arg CreateTopicParams) (Topic, error)
	GetTopic(id string) (Topic, error)
	ListTopics(arg ListTopicsParams) ([]Topic, error)

	CreateSubtopic(arg CreateSubtopicParams) (Subtopic, error)
	GetSubtopic(id string) (Subtopic, error)
	ListSubtopics(arg ListSubtopicsParams) ([]Subtopic, error)
}

type CreateSubjectParams struct {
	Name     string `json:"name" binding:"required"`
	Language string `json:"language" binding:"required,language"`
}

func subjectFromParams(arg CreateSubjectParams) Subject {
	return Subject{
		Name:     arg.Name,
		Language: arg.Language,
	}
}

func (db *MongoDB) CreateSubject(arg CreateSubjectParams) (Subject, error) {
	subject := subjectFromParams(arg)

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subjects")
	result, err := collection.InsertOne(context.Background(), subject)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	subject.ID = id

	return subject, err
}

func (db *MongoDB) GetSubject(id string) (Subject, error) {
	var subject Subject
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return subject, err
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subjects")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&subject)

	return subject, err
}

func (db *MongoDB) ListSubjects() ([]Subject, error) {
	var subjects []Subject

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subjects")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return subjects, err
	}

	err = cursor.All(context.Background(), &subjects)
	return subjects, err
}

type CreateTopicParams struct {
	Name      string `json:"name" binding:"required"`
	SubjectID string `json:"subject_id" binding:"required"`
}

func topicFromParams(arg CreateTopicParams) Topic {
	return Topic{
		Name:      arg.Name,
		SubjectID: arg.SubjectID,
	}
}

func (db *MongoDB) CreateTopic(arg CreateTopicParams) (Topic, error) {
	topic := topicFromParams(arg)

	if !db.isUniqueTopic(topic.Name) {
		return topic, errors.New("topic already exists")
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("topics")
	result, err := collection.InsertOne(context.Background(), topic)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	topic.ID = id

	return topic, err
}

func (db *MongoDB) isUniqueTopic(topic string) bool {
	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("topics")
	filter := bson.D{{Key: "name", Value: topic}}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	return count == 0
}

func (db *MongoDB) GetTopic(id string) (Topic, error) {
	var topic Topic
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return topic, err
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("topics")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&topic)

	return topic, err
}

type ListTopicsParams struct {
	SubjectID string `json:"subject_id" binding:"required"`
}

func (db *MongoDB) ListTopics(arg ListTopicsParams) ([]Topic, error) {
	var topics []Topic

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("topics")
	filter := bson.D{{Key: "subject_id", Value: arg.SubjectID}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return topics, err
	}

	err = cursor.All(context.Background(), &topics)
	return topics, err
}

type CreateSubtopicParams struct {
	Name    string `json:"name" binding:"required"`
	TopicID string `json:"topic_id" binding:"required"`
}

func subtopicFromParams(arg CreateSubtopicParams) Subtopic {
	return Subtopic{
		Name:    arg.Name,
		TopicID: arg.TopicID,
	}
}

func (db *MongoDB) CreateSubtopic(arg CreateSubtopicParams) (Subtopic, error) {
	subtopic := subtopicFromParams(arg)

	if !db.isUniqueSubtopic(subtopic.Name) {
		return subtopic, errors.New("subtopic already exists")
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subtopics")
	result, err := collection.InsertOne(context.Background(), subtopic)

	id := result.InsertedID.(primitive.ObjectID).Hex()
	subtopic.ID = id

	return subtopic, err
}

func (db *MongoDB) isUniqueSubtopic(subtopic string) bool {
	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subtopics")
	filter := bson.D{{Key: "name", Value: subtopic}}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}
	return count == 0
}

func (db *MongoDB) GetSubtopic(id string) (Subtopic, error) {
	var subtopic Subtopic
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return subtopic, err
	}

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subtopics")
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = collection.FindOne(context.Background(), filter).Decode(&subtopic)

	return subtopic, err
}

type ListSubtopicsParams struct {
	TopicID string `json:"topic_id" binding:"required"`
}

func (db *MongoDB) ListSubtopics(arg ListSubtopicsParams) ([]Subtopic, error) {
	var subtopics []Subtopic

	collection := db.client.Database(os.Getenv("MONGODB_DB_NAME")).Collection("subtopics")
	filter := bson.D{{Key: "topic_id", Value: arg.TopicID}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return subtopics, err
	}

	err = cursor.All(context.Background(), &subtopics)
	return subtopics, err
}
