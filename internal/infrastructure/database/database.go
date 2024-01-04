package database

import (
	"context"
	"kazokku/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, conf utils.DB) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, conf.DSN())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
