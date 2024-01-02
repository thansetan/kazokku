package database

import (
	"context"
	"fmt"
	"kazokku/internal/utils"

	"github.com/jackc/pgx/v5"
)

func New(ctx context.Context, conf utils.DB) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.UserName, conf.Password, conf.Host, conf.Port, conf.Name)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
