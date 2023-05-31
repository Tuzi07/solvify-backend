package db

type ProblemDatabase interface {
	ListProblems(arg ListProblemsParams) ([]Problem, error)
}

type ListProblemsParams struct {
}

func (dbManager *MongoDB) ListProblems(arg ListProblemsParams) ([]Problem, error) {
	return nil, nil
}
