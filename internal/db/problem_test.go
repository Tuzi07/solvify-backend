package db

import (
	"context"
	"testing"
	"github.com/Tuzi07/solvify-backend/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomProlemParams(t *testing.T) CreateProblemParams {
	arg := CreateProblemParams{
		Statement: util.RandomStatement(),
		Feedback: util.RandomFeedback(),
		SubjectId: util.RandomSubjectId(),
		TopicId: util.RandomTopicId(),
		SubtopicId: util.RandomSubtopicId(),
		LevelOfEducation: util.RandomLevelOfEducation(),
		Language: util.RandomLanguage(),
		CreatorID: util.RandomCreatorID(),
	}
	return arg
}

func createRandomTFProblem(t *testing.T) CreateTFProblemParams{
	arg := CreateTFProblemParams{
		CreateProblemParams: createRandomProlemParams(t),
		BoolAnswer:  util.RandomBoolAnswer(),
	}
	return arg
}

func createRandomMTFProblem(t *testing.T) CreateMTFProblemParams{
	items := util.RandomItems()
	arg := CreateMTFProblemParams{
		CreateProblemParams: createRandomProlemParams(t),
		Items: items,
		BoolAnswers: util.RandomBoolAnswers(len(items)),
	}
	return arg
}

func createRandomMCProblem(t *testing.T) CreateMCProblemParams{
	items := util.RandomItems()
	arg := CreateMCProblemParams{
		CreateProblemParams: createRandomProlemParams(t),
		Items: items,
		CorrectItem: util.RandomCorrectItem(len(items)),
	}
	return arg
}

func TestCreateTFProblem(t *testing.T) {
	arg := createRandomTFProblem(t)
	problem, err := testServer.CreateTFProblem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, problem)
	require.Equal(t, arg.Statement, problem.Statement)
	require.Equal(t, arg.Feedback, problem.Feedback)
	require.Equal(t, arg.SubjectId, problem.SubjectId)
	require.Equal(t, arg.TopicId, problem.TopicId)
	require.Equal(t, arg.SubtopicId, problem.SubtopicId)
	require.Equal(t, arg.LevelOfEducation, problem.LevelOfEducation)
	require.Equal(t, arg.Language, problem.Language)
	require.Equal(t, arg.CreatorID, problem.CreatorID)
	require.Equal(t, arg.BoolAnswer, problem.BoolAnswer)
	require.NotZero(t, problem.ID)
	require.NotZero(t, problem.CreatedAt)
}


