package task

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(task *Task) error
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
		return fmt.Errorf("could not create task: %w", err)
	}
	return nil
}
