package database

import (
	"context"
	"time"

	"github.com/Seven11Eleven/time-tracker-test-task/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct{
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
    return &TaskRepository{db: db}
}

func (r *TaskRepository) GetUserTasksByPeriod(ctx context.Context, userID int, start, end time.Time) ([]models.Task, error){
	var tasks []models.Task
	query := `
	SELECT id, user_id, description, start_time, end_time, created_at, updated_at
	FROM tasks
	WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
	ORDER BY EXTRACT(EPOCH FROM (end_time - start_time)) DESC
	`
	rows, err := r.db.Query(ctx, query, userID, start, end)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Description, &task.StartTime, &task.EndTime, &task.CreatedAt, &task.UpdatedAt); err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }
	if rows.Err() != nil{
		return nil, rows.Err()
	}

	return tasks, nil
	}

	func (r *TaskRepository) StartTask(ctx context.Context, userID int, description string) error {
		query := `
			INSERT INTO tasks (user_id, description, start_time, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW(), NOW())
		`
		_, err := r.db.Exec(ctx, query, userID, description)
		return err
	}
	
	func (r *TaskRepository) EndTask(ctx context.Context, taskID int) error {
		query := `
			UPDATE tasks
			SET end_time = NOW(), updated_at = NOW()
			WHERE id = $1 AND end_time IS NULL
		`
		_, err := r.db.Exec(ctx, query, taskID)
		return err
	}
	