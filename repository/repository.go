package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	client *pgxpool.Pool
}

func New() *Repository {
	r := new(Repository)
	//todo: set correct context
	pgxpool.Connect(context.Background(), "")
	return r
}