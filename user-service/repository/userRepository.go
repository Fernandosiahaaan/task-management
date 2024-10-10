package repository

import (
	"context"
	"database/sql"
	"user-service/internal/model"
)

type UserRepository struct {
	db     *sql.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func NewuserRepository(db *sql.DB, ctx context.Context) *UserRepository {
	dbCtx, dbCancel := context.WithCancel(ctx)
	return &UserRepository{
		db:     db,
		ctx:    dbCtx,
		cancel: dbCancel,
	}
}

func (r *UserRepository) CreateNewUser(user model.User) (string, error) {
	var id string
	query := `
	INSERT INTO users (id, username, password, email, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id 
	`
	err := r.db.QueryRowContext(r.ctx, query, user.Id, user.Username, user.Password, user.Email, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&id)
	return id, err
}

func (r *UserRepository) UpdateUser(user model.User) (string, error) {
	var id string
	query := `
        UPDATE users 
        SET username = $1, password = $2, email = $3, role = $4, updated_at = $5
        WHERE id = $6
        RETURNING id
    `
	err := r.db.QueryRowContext(r.ctx, query, user.Username, user.Password, user.Email, user.Role, user.UpdatedAt, user.Id).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserRepository) GetUser(user model.User) (model.User, error) {
	query := `
	SELECT id, username, password, email, role, created_at, updated_at 
	FROM users 
	WHERE username=$1
	`
	var existUser model.User
	err := r.db.QueryRowContext(r.ctx, query, user.Username).Scan(
		&existUser.Id,
		&existUser.Username,
		&existUser.Password,
		&existUser.Email,
		&existUser.Role,
		&existUser.CreatedAt,
		&existUser.UpdatedAt,
	)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (r *UserRepository) GetUserById(userId string) (*model.User, error) {
	query := `
	SELECT id, username, password, email, role, created_at, updated_at
	FROM users 
	WHERE id=$1
	`
	var existUser *model.User = &model.User{}
	err := r.db.QueryRowContext(r.ctx, query, userId).Scan(
		&existUser.Id,
		&existUser.Username,
		&existUser.Password,
		&existUser.Email,
		&existUser.Role,
		&existUser.CreatedAt,
		&existUser.UpdatedAt,
	)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	query := `
	SELECT id, username, password, email, role, created_at, updated_at
	FROM users
	`
	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Close() {
	r.db.Close()
	r.cancel()
}
