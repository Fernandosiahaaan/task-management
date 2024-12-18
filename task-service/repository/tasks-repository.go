package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"task-service/internal/model"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type TaskRepository struct {
	DB     *sql.DB
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewTaskRepository(ctx context.Context) (*TaskRepository, error) {
	sqltrace.Register("postgres", &pq.Driver{})

	db, err := sqltrace.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	repoCtx, repoCancel := context.WithCancel(ctx)
	return &TaskRepository{
		DB:     db,
		Ctx:    repoCtx,
		Cancel: repoCancel,
	}, nil
}

func (r *TaskRepository) CreateNewTask(task *model.Task) (int64, error) {
	var id int64
	query := `
    INSERT INTO tasks (title, description, status, due_date, assigned_to, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id 
    `
	// Perhatikan bahwa id sekarang berupa nilai (int64), bukan pointer
	err := r.DB.QueryRowContext(r.Ctx, query,
		task.Title,
		task.Description,
		task.Status,
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
	SELECT id, title, description, status, due_date, assigned_to, created_by, updated_by, created_at, updated_at
	FROM tasks
	WHERE id=$1 
	`

	err := r.DB.QueryRowContext(r.Ctx, query, taskId).Scan(
		&returntask.Id,
		&returntask.Title,
		&returntask.Description,
		&returntask.Status,
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
	SELECT id, title, description, status, due_date, assigned_to, created_by, updated_by, created_at, updated_at
	FROM tasks
	`
	// Execute query
	rows, err := r.DB.QueryContext(r.Ctx, query)
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
			&task.Status,
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
	SET title = $1, description = $2, status = $3, due_date = $4, assigned_to = $5, updated_by = $6, updated_at = $7
	WHERE id = $8 
	RETURNING id 
	`
	err := r.DB.QueryRowContext(r.Ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.AssignedTo,
		task.UpdatedBy,
		task.UpdatedAt,
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
	res, err := r.DB.ExecContext(r.Ctx, query, taskId)
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

func (r *TaskRepository) GetAllTaskByUserId(userID string) ([]*model.Task, error) {
	query := `
	SELECT id, title, description, status, due_date, assigned_to, created_by, updated_by, created_at, updated_at
	FROM tasks
	WHERE assigned_to=$1
	`
	// Execute query
	rows, err := r.DB.QueryContext(r.Ctx, query, userID)
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
			&task.Status,
			&task.DueDate,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.UpdatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Close() {
	r.DB.Close()
	r.Cancel()
}
