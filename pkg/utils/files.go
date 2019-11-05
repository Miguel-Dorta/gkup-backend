package utils

import (
	"fmt"
	"io"
	"os"
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
