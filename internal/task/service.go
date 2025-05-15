package task

import "github.com/google/uuid"

type Service interface {
	CreateTask(userID uuid.UUID, task *TaskRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}
func (s *service) CreateTask(userID uuid.UUID, task *TaskRequest) error {
	//if err := s.repo.Create(task); err != nil {
	//	return nil, fmt.Errorf("could not create task: %w", err)
	//}

	//t := &Task{
	//	UserID: userID,
	//	Title:  task.Title,
	//	Description: task.Description,
	//	Status: "pending",
	//	Priority: task.Priority,
	//}
	return s.repo.Create(&Task{
		UserID:      userID,
		Title:       task.Title,
		Description: task.Description,
		Status:      "pending",
		Priority:    task.Priority,
	})
}
