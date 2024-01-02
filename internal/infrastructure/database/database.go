package database

import (
	"context"
	"fmt"
	"kazokku/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, conf utils.DB) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.UserName, conf.Password, conf.Host, conf.Port, conf.Name)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
