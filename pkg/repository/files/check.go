package files

import (
	"bytes"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"os"
	"path/filepath"
)

func Check(path string, h *hash.Hasher) error {
	// Get expected values
	expectedHash, expectedSize, err := GetDataFromName(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("error parsing filename (%s): %w", path, err)
	}

	// Check size
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot get info from file (%s): %w", path, err)
	}
	if stat.Size() != expectedSize {
		return fmt.Errorf("size found (%d) doesn't match expected size (%d)", stat.Size(), expectedSize)
	}

	// Check hash
	actualHash, err := h.HashFile(path)
	if err != nil {
		return fmt.Errorf("error hashing file (%s): %w", path, err)
	}
	if !bytes.Equal(expectedHash, actualHash) {
		return fmt.Errorf("hash found (%x) doesn't match expected hash (%x)", actualHash, expectedHash)
	}

	return nil
}
