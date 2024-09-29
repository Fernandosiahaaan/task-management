package repository

import (
	"context"
	"database/sql"
	"task-management/user-service/internal/model"
)

type UserRepository struct {
	DB  *sql.DB
	Ctx context.Context
}

func NewuserRepository(db *sql.DB, ctx context.Context) *UserRepository {
	return &UserRepository{DB: db, Ctx: ctx}
}

func (r *UserRepository) CreateNewUser(user model.User) (string, error) {
	var id string
	query := `
	INSERT INTO users (id, username, password, email, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id 
	`
	err := r.DB.QueryRowContext(r.Ctx, query, user.Id, user.Username, user.Password, user.Email, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&id)
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
	err := r.DB.QueryRowContext(r.Ctx, query, user.Username, user.Password, user.Email, user.Role, user.UpdatedAt, user.Id).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserRepository) GetUser(user model.User) (model.User, error) {
	query := `
	SELECT id, username, password, email, role FROM users 
	WHERE username=$1
	`
	var existUser model.User
	err := r.DB.QueryRowContext(r.Ctx, query, user.Username).Scan(
		&existUser.Id,
		&existUser.Username,
		&existUser.Password,
		&existUser.Email,
		&existUser.Role,
	)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	query := `
	SELECT id, username, role 
	FROM users
	`
	rows, err := r.DB.QueryContext(r.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Username, &user.Role)
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

// func (r *UserRepository) GetUserById(id int64) (int64, error) {
// 	var id int64
// 	query := `
// 	INSERT INTO users (username, pasword_hash, email)
// 	VALUES ($1, $2, $3) RETURNING id
// 	`
// 	err := r.DB.QueryRowContext(r.Ctx, query, user.Username, user.Password, user.Email).Scan(&id)
// 	return id, err
// }
