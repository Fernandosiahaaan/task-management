package repository

import (
	"context"
	"database/sql"
	"task-management/task-service/internal/model"
)

type TaskRepository struct {
	DB  *sql.DB
	Ctx context.Context
}

func NewTaskRepository(db *sql.DB, ctx context.Context) *TaskRepository {
	return &TaskRepository{DB: db, Ctx: ctx}
}

func (r *TaskRepository) CreateNewTask(task *model.Task) (*int64, error) {
	var id *int64
	query := `
	INSERT INTO task (id, title, description, due_date, assigned_to, created_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id 
	`
	err := r.DB.QueryRowContext(r.Ctx, query, task.Id, task.Title,
		task.Description, task.Description, task.DueDate, task.AssignedTo,
		task.CreatedBy, task.CreatedAt, task.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *TaskRepository) GetTask(task *model.Task) (*model.Task, error) {
	var returntask model.Task
	query := `
	SELECT id, title, description, due_date, assigned_to, created_by
	WHERE id=$1 
	`

	err := r.DB.QueryRowContext(r.Ctx, query, task.Id).Scan(
		&returntask.Id,
		&returntask.Title,
		&returntask.Description,
		&returntask.DueDate,
		&returntask.AssignedTo,
		&returntask.AssignedTo,
		&returntask.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &returntask, nil
}

func (r *TaskRepository) UpdateTask(task *model.Task) (*int64, error) {
	var id *int64
	query := `
	UPDATE task 
	SET title = $1, description = $2, due_date = $3, assigned_to = $4, created_by = $5, created_at = $6, updated_at = $7
	WHERE id = $8 
	RETURNING id 
	`
	err := r.DB.QueryRowContext(r.Ctx, query, task.Title, task.Description,
		task.Description, task.DueDate, task.AssignedTo,
		task.CreatedBy, task.CreatedAt, task.UpdatedAt,
		task.Id,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return id, nil
}
