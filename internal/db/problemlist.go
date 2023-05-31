package db

// type ProblemListDatabase interface {
// 	CreateProblemList(list *models.ProblemList) error
// 	GetProblemList(list *models.ProblemList, id string) error
// 	UpdateProblemList(list *models.ProblemList, id string) error
// 	DeleteProblemList(id string) error
// }

// func (dbManager *MongoDB) CreateProblemList(list *models.ProblemList) error {
// 	list.CreatedAt = time.Now()
// 	id, err := dbManager.create(&list, "problemlists")

// 	list.ID = id
// 	return err
// }

// func (dbManager *MongoDB) GetProblemList(list *models.ProblemList, id string) error {
// 	err := dbManager.get(&list, id, "problemlists")
// 	return err
// }

// func (dbManager *MongoDB) UpdateProblemList(list *models.ProblemList, id string) error {
// 	err := dbManager.update(&list, id, "problemlists")

// 	list.ID = id
// 	return err
// }

// func (dbManager *MongoDB) DeleteProblemList(id string) error {
// 	err := dbManager.delete(id, "problemlists")
// 	return err
// }
