package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(userID uint, file *multipart.FileHeader) (string, error) {
	outputDir := filepath.Join(os.Getenv("SAVE_DIR"), "photos", fmt.Sprintf("%d", userID))
	outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename)))

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return outputFilePath, err
	}

	dst, err := os.Create(outputFilePath)
	if err != nil {
		return outputFilePath, err
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return outputFilePath, err
	}
	defer src.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return outputFilePath, err
	}

	return filepath.ToSlash(strings.TrimPrefix(outputFilePath, filepath.Join(os.Getenv("SAVE_DIR"), "photos"))), nil
}

func IsImage(fileType string) bool {
	return strings.HasPrefix(fileType, "image/")
}
