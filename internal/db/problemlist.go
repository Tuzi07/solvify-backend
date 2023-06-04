package db

// type ProblemListDatabase interface {
// 	CreateProblemList(list *models.ProblemList) error
// 	GetProblemList(list *models.ProblemList, id string) error
// 	UpdateProblemList(list *models.ProblemList, id string) error
// 	DeleteProblemList(id string) error
// }

// func (db *MongoDB) CreateProblemList(list *models.ProblemList) error {
// 	list.CreatedAt = time.Now()
// 	id, err := db.create(&list, "problemlists")

// 	list.ID = id
// 	return err
// }

// func (db *MongoDB) GetProblemList(list *models.ProblemList, id string) error {
// 	err := db.get(&list, id, "problemlists")
// 	return err
// }

// func (db *MongoDB) UpdateProblemList(list *models.ProblemList, id string) error {
// 	err := db.update(&list, id, "problemlists")

// 	list.ID = id
// 	return err
// }

// func (db *MongoDB) DeleteProblemList(id string) error {
// 	err := db.delete(id, "problemlists")
// 	return err
// }
