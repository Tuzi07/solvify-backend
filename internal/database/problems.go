package database

import (
	"github.com/Tuzi07/solvify-backend/internal/models"
)

type ProblemDatabase interface {
	CreateProblem(problem *models.Problem) error
	GetProblem(problem *models.Problem, id string) error
	UpdateProblem(problem *models.Problem, id string) error
	DeleteProblem(id string) error
}

func (dbManager *MongoDB) CreateProblem(problem *models.Problem) error {
	id, err := dbManager.create(&problem, "problems")

	problem.ID = id
	return err
}

func (dbManager *MongoDB) GetProblem(problem *models.Problem, id string) error {
	err := dbManager.get(&problem, id, "problems")
	return err
}

func (dbManager *MongoDB) UpdateProblem(problem *models.Problem, id string) error {
	err := dbManager.update(&problem, id, "problems")

	problem.ID = id
	return err
}

func (dbManager *MongoDB) DeleteProblem(id string) error {
	err := dbManager.delete(id, "problems")
	return err
}
