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
