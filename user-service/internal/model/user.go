package model

import "time"

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password_hash"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
	ID int64 `json:"id"`
}
