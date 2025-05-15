package task

import "github.com/google/uuid"

type Service interface {
	CreateTask(userID uuid.UUID, task *TaskRequest) error
	ListTasks(userID uuid.UUID) ([]*Task, error)
	GetTask(userID uuid.UUID, taskID uuid.UUID) (*Task, error)
	DeleteTask(userID uuid.UUID, taskID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}
func (s *service) CreateTask(userID uuid.UUID, task *TaskRequest) error {
	return s.repo.Create(&Task{
		UserID:      userID,
		Title:       task.Title,
		Description: task.Description,
		Status:      "pending",
		Priority:    task.Priority,
	})
}

func (s *service) ListTasks(userID uuid.UUID) ([]*Task, error) {

	tasks, err := s.repo.List(userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func (s *service) GetTask(userID uuid.UUID, taskID uuid.UUID) (*Task, error) {

	task, err := s.repo.GetTask(userID, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *service) DeleteTask(userID uuid.UUID, taskID uuid.UUID) error {
	return s.repo.DeleteTask(userID, taskID)
}
