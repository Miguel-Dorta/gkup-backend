package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from origin path to destiny path
func CopyFile(origin, destiny string, buffer []byte) error {
	originFile, err := os.Open(origin)
	if err != nil {
		return fmt.Errorf("cannot open file (%s): %w", origin, err)
	}
	defer originFile.Close()

	destinyFile, err := os.Create(destiny)
	if err != nil {
		return fmt.Errorf("cannot create file (%s): %w", destiny, err)
	}
	defer destinyFile.Close()

	if _, err = io.CopyBuffer(destinyFile, originFile, buffer); err != nil {
		return fmt.Errorf("error copying file from \"%s\" to \"%s\": %w", origin, destiny, err)
	}

	if err = destinyFile.Close(); err != nil {
		return fmt.Errorf("error closing file (%s): %w", destiny, err)
	}
	return nil
}

// ListDir lists the directory from the path provided
func ListDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdir(-1)
}

// CreateWithParents is similar to os.Create() but creating the parent directories necessaries.
func CreateWithParents(path string) (*os.File, error) {
	parentPath := filepath.Dir(path)

	parentExists, err := FileExist(parentPath)
	if err != nil {
		return nil, fmt.Errorf("error checking parent dir existence of file \"%s\": %w", path, err)
	}

	if !parentExists {
		if err = os.MkdirAll(parentPath, 0755); err != nil {
			return nil, fmt.Errorf("error creating parent dirs of file \"%s\": %w", path, err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("error creating file \"%s\": %w", path, err)
	}
	return f, nil
}

// FileExist return whether a file exist or not. It returns an error if it cannot determine it.
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
