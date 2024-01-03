package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func NewLogger(logOutputDir string) (*slog.Logger, error) {
	err := os.MkdirAll(logOutputDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(filepath.Join(logOutputDir, "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stderr, logFile)
	logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{
		AddSource: true,
	}))

	return logger, nil
}
