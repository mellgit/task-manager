package task

import (
	"github.com/google/uuid"
	"github.com/mellgit/task-manager/internal/queue"
)

type Service interface {
	CreateTask(userID uuid.UUID, task *TaskRequest) error
	ListTasks(userID uuid.UUID) ([]*Task, error)
	GetTask(userID uuid.UUID, taskID uuid.UUID) (*Task, error)
	DeleteTask(userID uuid.UUID, taskID uuid.UUID) error
	UpdateTask(userID uuid.UUID, task *TaskRequest) error
}

type service struct {
	repo     Repository
	producer *queue.Producer
}

func NewService(repo Repository, producer *queue.Producer) Service {
	return &service{repo, producer}
}
func (s *service) CreateTask(userID uuid.UUID, task *TaskRequest) error {

	taskID, err := s.repo.Create(&Task{
		UserID:      userID,
		Title:       task.Title,
		Description: task.Description,
		Status:      "pending",
		Priority:    task.Priority,
	})
	if err != nil {
		return err
	}
	if err := s.producer.Publish(queue.TaskPayload{TaskID: uuid.MustParse(taskID), UserID: userID}); err != nil {
		return err
	}
	return nil

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

func (s *service) UpdateTask(userID uuid.UUID, task *TaskRequest) error {

	return s.repo.UpdateTask(userID, &Task{
		UserID:      userID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
	})
}
