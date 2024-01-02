package repository

import (
	"context"
	"kazokku/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PhotoRepository interface {
	InsertBatch(context.Context, pgx.Tx, []domain.Photo) error
}

type photoRepository struct {
	db *pgx.Conn
}

func NewPhotoRepository(db *pgx.Conn) photoRepository {
	return photoRepository{db}
}

func (repo photoRepository) InsertBatch(ctx context.Context, tx pgx.Tx, data []domain.Photo) error {
	stmt := "INSERT INTO photos(user_id, filename) VALUES ($1, $2);"
	batch := new(pgx.Batch)

	for _, photo := range data {
		batch.Queue(stmt, photo.UserID, photo.Filepath)
	}

	res := tx.SendBatch(ctx, batch)
	defer res.Close()
	if _, err := res.Exec(); err != nil {
		return err
	}

	return nil
}
