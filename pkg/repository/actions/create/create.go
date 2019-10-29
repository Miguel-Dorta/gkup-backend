package create

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"path/filepath"
)

func Create(path, hashAlgorithm string) error {
	// Get path stat
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) { // If it's not a "not exist" error, return it
			return &os.PathError{
				Op:   "stat repository path",
				Path: path,
				Err:  err,
			}
		}
		return create(path, hashAlgorithm) // If it's a "not exist" error, create it. Done.
	}

	// Check if it's not a dir
	if !stat.IsDir() {
		return errors.New("repository path must be a directory")
	}

	// Check if dir is empty
	list, err := utils.ListDir(path)
	if err != nil {
		return &os.PathError{
			Op:   "list repository directory",
			Path: path,
			Err:  err,
		}
	}
	if len(list) != 0 {
		return errors.New("repository path must be empty")
	}

	return create(path, hashAlgorithm)
}

// create creates a repository in the path provided with the algorithm provided.
// the path must exist and be an empty directory.
func create(path, hashAlgorithm string) error {
	// Create snapshots dir
	snapshotsFolderPath := filepath.Join(path, repository.SnapshotsFolderName)
	if err := os.MkdirAll(snapshotsFolderPath, pkg.DefaultDirPerm); err != nil {
		return &os.PathError{
			Op:   "create snapshots folder",
			Path: snapshotsFolderPath,
			Err:  err,
		}
	}

	// Create files dir and subdirectories
	filesFolderPath := filepath.Join(path, repository.FilesFolderName)
	for i:=0; i<=0xff; i++ {
		subDirPath := filepath.Join(filesFolderPath, fmt.Sprintf("%02x", i))
		if err := os.MkdirAll(subDirPath, pkg.DefaultDirPerm); err != nil {
			return &os.PathError{
				Op:   "create files folders",
				Path: subDirPath,
				Err:  err,
			}
		}
	}

	// Create settings file
	if err := settings.Write(filepath.Join(path, settings.FileName), hashAlgorithm); err != nil {
		return fmt.Errorf("error creating settings file: %s", err)
	}
	return nil
}
