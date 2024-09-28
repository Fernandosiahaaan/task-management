package repository

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	DB  *sql.DB
	Ctx context.Context
}

func NewuserRepository(db *sql.DB, ctx context.Context) *UserRepository {
	return &UserRepository{DB: db, Ctx: ctx}
}
