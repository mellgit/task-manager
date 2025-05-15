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

	/*
		INSERT INTO tasks (
		    user_id,
		    title,
		    description,
		    status,
		    priority,
		    created_at,
		    updated_at
		) VALUES (
		    '37a90e48-a257-46f4-8522-553fe991bb36',
		    'Изучить Golang',
		    'Пройти курс по Golang и написать тестовое приложение',
		    'in_progress',
		    2,
		    CURRENT_TIMESTAMP,
		    CURRENT_TIMESTAMP
		);
	*/

	query := `insert into tasks (user_id, title, description, status, priority, created_at, updated_at)  
				values ($1, $2, $3, $4, $5, NOW(), NOW())`

	//_ = r.db.QueryRowContext(r.ctx, query, task.UserID, task.Title, task.Description, task.Status, task.Priority)
	_, err := r.db.ExecContext(r.ctx, query, task.UserID, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}
	return nil
}
