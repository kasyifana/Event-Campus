package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file size exceeds maximum limit")
	ErrNoFile          = errors.New("no file provided")
)

// FileUploader handles file uploads
type FileUploader struct {
	uploadPath string
	maxSize    int64
}

// NewFileUploader creates a new file uploader
func NewFileUploader(uploadPath string, maxSize int64) *FileUploader {
	return &FileUploader{
		uploadPath: uploadPath,
		maxSize:    maxSize,
	}
}

// SavePoster saves an event poster
func (u *FileUploader) SavePoster(file *multipart.FileHeader) (string, error) {
	// Validate file size
	if file.Size > u.maxSize {
		return "", ErrFileTooLarge
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", ErrInvalidFileType
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	posterPath := filepath.Join(u.uploadPath, "posters", filename)

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(posterPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(posterPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path
	return filepath.Join("posters", filename), nil
}

// SaveDocument saves a whitelist document (PDF)
func (u *FileUploader) SaveDocument(file *multipart.FileHeader) (string, error) {
	// Validate file size
	if file.Size > u.maxSize {
		return "", ErrFileTooLarge
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return "", ErrInvalidFileType
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	docPath := filepath.Join(u.uploadPath, "documents", filename)

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(docPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Save file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(docPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path
	return filepath.Join("documents", filename), nil
}

// DeleteFile deletes a file from storage
func (u *FileUploader) DeleteFile(relativePath string) error {
	if relativePath == "" {
		return nil
	}

	fullPath := filepath.Join(u.uploadPath, relativePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFilePath returns the full path of a file
func (u *FileUploader) GetFilePath(relativePath string) string {
	return filepath.Join(u.uploadPath, relativePath)
}
