package main

import (
	"context"
	"kazokku/internal/infrastructure/database"
	"kazokku/internal/infrastructure/http"
	"kazokku/internal/utils"
	"log/slog"
	"os"
)

func main() {
	conf, err := utils.LoadConfig(".env")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger, err := utils.NewLogger(conf.App.SaveDir)
	if err != nil {
		slog.Error("failed to create logger instance", "error", err)
		os.Exit(1)
	}

	err = database.Migrate(conf.Database)
	if err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := database.New(ctx, conf.Database)
	if err != nil {
		logger.Error("failed to create database instance", "error", err)
		os.Exit(1)
	}

	app := http.New(conf.App, db, logger)
	if err := app.Run(); err != nil {
		logger.Error("failed to start app", "error", err)
		os.Exit(1)
	}
}
