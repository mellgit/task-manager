package worker

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Service interface {
	Process(payload TaskPayload) error
	GetPayload() TaskPayload
}

type service struct {
	repo Repository
	log  *log.Entry
}

func NewService(repo Repository, log *log.Entry) Service {
	//func NewService(repo Repository, consumer *queue.Consumer) Service {
	return &service{repo, log}
}

func (s *service) Process(payload TaskPayload) error {

	// update the issue status to "in_progress"
	if err := s.repo.ChangeStatus("in_progress", payload.TaskID); err != nil {
		return err
	}
	log.Infof("processing task %v", payload.TaskID)

	// complete the task
	s.doSomething()

	// update the issue status to "done" or "failed"
	if err := s.repo.ChangeStatus("done", payload.TaskID); err != nil {
		return err
	}
	log.Infof("task is done %v", payload.TaskID)

	return nil
}

func (s *service) GetPayload() TaskPayload {
	return TaskPayload{}
}

func (s *service) doSomething() {
	// I did not create a separate functionality for the worker, as it can be any, depending on the task.
	time.Sleep(5 * time.Second)
}
