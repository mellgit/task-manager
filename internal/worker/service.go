package worker

import (
	"fmt"
	"github.com/mellgit/task-manager/internal/queue"
	"time"
)

type Service interface {
	Process(payload TaskPayload) error
}

type service struct {
	repo     Repository
	consumer *queue.Consumer
}

func NewService(repo Repository, consumer *queue.Consumer) Service {
	return &service{repo, consumer}
}

func (s *service) Process(payload TaskPayload) error {

	// update the issue status to "in_progress"
	if err := s.repo.ChangeStatus("in_progress", payload.TaskID); err != nil {
		return err
	}
	fmt.Println("processing task", payload.TaskID)

	// complete the task
	s.doSomething()

	// update the issue status to "done" or "failed"
	if err := s.repo.ChangeStatus("done", payload.TaskID); err != nil {
		return err
	}
	fmt.Println("done task", payload.TaskID)

	return nil
}

func (s *service) doSomething() {
	// I did not create a separate functionality for the worker, as it can be any, depending on the task.
	time.Sleep(5 * time.Second)
}
