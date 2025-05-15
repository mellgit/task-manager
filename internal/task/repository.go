package task

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type Repository interface {
	Create(task *Task) error
	List(userID uuid.UUID) ([]*Task, error)
	GetTask(userID uuid.UUID, taskID uuid.UUID) (*Task, error)
	DeleteTask(userID uuid.UUID, taskID uuid.UUID) error
	UpdateTask(userID uuid.UUID, task *Task) error
	/*
		todo
		get task
		delete task
		get all tasks
		update task
	*/
}

type repository struct {
	ctx context.Context
	db  *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repository{
		ctx: context.Background(),
		db:  db,
	}
}

func (r *repository) Create(task *Task) error {

	query := `insert into tasks (user_id, title, description, status, priority, created_at, updated_at)  
				values ($1, $2, $3, $4, $5, NOW(), NOW())`

	_, err := r.db.ExecContext(r.ctx, query, task.UserID, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return fmt.Errorf("exec ctx: %w", err)
	}
	return nil
}

func (r *repository) List(userID uuid.UUID) ([]*Task, error) {
	query := `select * from tasks where user_id = $1`
	rows, err := r.db.QueryContext(r.ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query ctx: %w", err)
	}
	defer rows.Close()
	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		tasks = append(tasks, &task)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return tasks, nil
}

func (r *repository) GetTask(userID uuid.UUID, taskID uuid.UUID) (*Task, error) {

	query := `select * from tasks where user_id = $1 and id = $2`
	row := r.db.QueryRowContext(r.ctx, query, userID, taskID)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	var task Task
	if err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return &task, nil

}

func (r *repository) DeleteTask(userID uuid.UUID, taskID uuid.UUID) error {
	query := `delete from tasks where user_id = $1 and id = $2`
	_, err := r.db.ExecContext(r.ctx, query, userID, taskID)
	if err != nil {
		return fmt.Errorf("exec ctx: %w", err)
	}
	return nil
}

func (r *repository) UpdateTask(userID uuid.UUID, task *Task) error {

	query := `update tasks set title=$1, description=$2, priority=$3, updated_at=NOW() where user_id = $4`
	_, err := r.db.ExecContext(r.ctx, query, task.Title, task.Description, task.Priority, userID)
	if err != nil {
		return fmt.Errorf("exec ctx: %w", err)
	}
	return nil

}
