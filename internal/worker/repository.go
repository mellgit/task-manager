package worker

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type Repository interface {
	ChangeStatus(status string, taskID uuid.UUID) error
}

type repository struct {
	ctx context.Context
	db  *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repository{ctx: context.Background(), db: db}
}
func (r *repository) ChangeStatus(status string, taskID uuid.UUID) error {
	query := `UPDATE tasks
	SET status = $1, updated_at = CURRENT_TIMESTAMP
	WHERE id = $2`

	_, err := r.db.Exec(query, status, taskID)
	if err != nil {
		return fmt.Errorf("exec ctx: %w", err)
	}
	return nil
}
