package files

import (
	"bytes"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"path/filepath"
)

func Add(filesDirPath, filePath string, f *File, copyBuf []byte) error {
	destination := filepath.Join(filesDirPath, getSavingPath(f.Hash, f.Size))

	exists, err := utils.FileExist(destination)
	if err != nil {
		return fmt.Errorf("error checking file existence: %s", err)
	}

	if exists {
		return nil
	}

	if err := utils.CopyFile(filePath, destination, copyBuf); err != nil {
		return fmt.Errorf("error copying file (%s) to repository: %s", filePath, err)
	}
	return nil
}

func Restore(filesPath, restorationPath string, f *File, copyBuf []byte) error {
	return utils.CopyFile(
		filepath.Join(filesPath, getSavingPath(f.Hash, f.Size)),
		filepath.Join(restorationPath, filepath.FromSlash(f.RelativePath)),
		copyBuf,
	)
}

func Check(path string, h *hash.Hasher) error {
	// Get expected values
	expectedHash, expectedSize, err := getDataFromName(filepath.Base(path))
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
