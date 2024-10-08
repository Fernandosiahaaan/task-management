package repository

import (
	"context"
	"database/sql"
	"fmt"
	"task-service/internal/model"
)

type TaskRepository struct {
	db     *sql.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTaskRepository(db *sql.DB, ctx context.Context) *TaskRepository {
	repoCtx, repoCancel := context.WithCancel(ctx)
	return &TaskRepository{
		db:     db,
		ctx:    repoCtx,
		cancel: repoCancel,
	}
}

func (r *TaskRepository) CreateNewTask(task *model.Task) (int64, error) {
	var id int64
	query := `
    INSERT INTO tasks (title, description, due_date, assigned_to, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id 
    `
	// Perhatikan bahwa id sekarang berupa nilai (int64), bukan pointer
	err := r.db.QueryRowContext(r.ctx, query,
		task.Title,
		task.Description,
		task.DueDate,
		task.AssignedTo,
		task.CreatedBy,
		task.UpdatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TaskRepository) GetTaskById(taskId *int64) (*model.Task, error) {
	var returntask model.Task
	query := `
	SELECT id, title, description, due_date, assigned_to, created_by, updated_by, created_at, updated_at
	FROM tasks
	WHERE id=$1 
	`

	err := r.db.QueryRowContext(r.ctx, query, taskId).Scan(
		&returntask.Id,
		&returntask.Title,
		&returntask.Description,
		&returntask.DueDate,
		&returntask.AssignedTo,
		&returntask.CreatedBy,
		&returntask.UpdatedBy,
		&returntask.CreatedAt,
		&returntask.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &returntask, nil
}
func (r *TaskRepository) GetAllTask() ([]*model.Task, error) {
	query := `
	SELECT id, title, description, due_date, assigned_to, created_by, updated_by, created_at, updated_at
	FROM tasks
	`
	// Execute query
	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.UpdatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get all tasks from db: %s", err.Error())
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(task *model.Task) (*int64, error) {
	var id *int64
	query := `
	UPDATE tasks 
	SET title = $1, description = $2, due_date = $3, assigned_to = $4, updated_by = $5, created_at = $6, updated_at = $7
	WHERE id = $8 
	RETURNING id 
	`
	err := r.db.QueryRowContext(r.ctx, query, task.Title, task.Description,
		task.DueDate, task.AssignedTo,
		task.UpdatedBy, task.CreatedAt, task.UpdatedAt,
		task.Id,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed update database. err = %s", err)
	}
	return id, nil
}

func (r *TaskRepository) DeleteTask(taskId *int64) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1
	`
	res, err := r.db.ExecContext(r.ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("failed delete task from database. err = %s", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed check affected of delete task from database. err = %s", err)
	}

	fmt.Printf("Task with ID %d deleted successfully. Rows affected: %d\n", taskId, rowAffected)
	return nil
}

func (r *TaskRepository) Close() {
	r.db.Close()
	r.cancel()
}
