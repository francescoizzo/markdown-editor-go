package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileUtils provides utilities for file operations
type FileUtils struct{}

// SaveToFile writes content to a file
func (f *FileUtils) SaveToFile(path string, content string) error {
	if path == "" {
		return errors.New("file path cannot be empty")
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// Write content to file
	return ioutil.WriteFile(path, []byte(content), 0644)
}

// ReadFromFile reads content from a file
func (f *FileUtils) ReadFromFile(path string) (string, error) {
	if path == "" {
		return "", errors.New("file path cannot be empty")
	}

	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", errors.New("file does not exist")
	}

	// Read file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetFilenameFromPath extracts the filename from a path
func (f *FileUtils) GetFilenameFromPath(path string) string {
	return filepath.Base(path)
}

// EnsureExtension ensures the file has a .md extension
func (f *FileUtils) EnsureExtension(path string) string {
	if path == "" {
		return "untitled.md"
	}

	// Check if path already has .md or .markdown extension
	lowerPath := strings.ToLower(path)
	if strings.HasSuffix(lowerPath, ".md") || strings.HasSuffix(lowerPath, ".markdown") {
		return path
	}

	// Add .md extension
	return path + ".md"
}

// CreateTempBackup creates a temporary backup of a file
func (f *FileUtils) CreateTempBackup(path string) (string, error) {
	if path == "" {
		return "", errors.New("file path cannot be empty")
	}

	// Get directory and filename
	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	// Create backup filename
	backupPath := filepath.Join(dir, "."+filename+".bak")

	// Read original file
	content, err := f.ReadFromFile(path)
	if err != nil {
		return "", err
	}

	// Write to backup file
	err = f.SaveToFile(backupPath, content)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

// IsFileModifiedExternally checks if file has been modified since last read
func (f *FileUtils) IsFileModifiedExternally(path string, lastModTime int64) (bool, error) {
	if path == "" {
		return false, errors.New("file path cannot be empty")
	}

	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// Compare modification times
	return fileInfo.ModTime().Unix() > lastModTime, nil
}
