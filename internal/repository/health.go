package repository

import (
	"context"
	"database/sql"
)

type HealthRepository interface {
	CheckDB(ctx context.Context) error
}

type healthRepository struct {
	db *sql.DB
}

func NewHealthRepository(db *sql.DB) HealthRepository {
	return &healthRepository{
		db: db,
	}
}

func (r *healthRepository) CheckDB(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
