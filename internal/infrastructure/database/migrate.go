package database

import (
	"database/sql"
	"kazokku/internal/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(conf utils.DB) error {
	db, err := sql.Open("pgx", conf.DSN())
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{
		DatabaseName: conf.Name,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		conf.Name,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
