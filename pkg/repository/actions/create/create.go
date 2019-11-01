package create

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	"path/filepath"
)

func Create(path, hashAlgorithm string, errWriter io.Writer) {
	// Get path stat
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) { // If it's not a "not exist" error, return it
			_, _ = fmt.Fprintf(errWriter, "error getting info from path \"%s\": %s\n", path, err)
			return
		}

		// If it's a "not exist" error, create it. Done.
		if err := create(path, hashAlgorithm); err != nil {
			_, _ = fmt.Fprintln(errWriter, err)
		}
		return // Finalization successfully
	}

	// Check if it's not a dir
	if !stat.IsDir() {
		_, _ = fmt.Fprintln(errWriter, "repository path must be a directory")
		return
	}

	// Check if dir is empty
	list, err := utils.ListDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "error listing repository directory (%s): %s\n", path, err)
		return
	}
	if len(list) != 0 {
		_, _ = fmt.Fprintln(errWriter, "repository path must be empty")
		return
	}

	if err := create(path, hashAlgorithm); err != nil {
		_, _ = fmt.Fprintln(errWriter, err)
	}
}

// create creates a repository in the path provided with the algorithm provided.
// the path must exist and be an empty directory.
func create(path, hashAlgorithm string) error {
	// Create snapshots dir
	snapshotsFolderPath := filepath.Join(path, snapshots.FolderName)
	if err := os.MkdirAll(snapshotsFolderPath, pkg.DefaultDirPerm); err != nil {
		return fmt.Errorf("error creating snapshot folder (%s): %s", snapshotsFolderPath, err)
	}

	// Create files dir and subdirectories
	filesFolderPath := filepath.Join(path, files.FolderName)
	for i:=0; i<=0xff; i++ {
		subDirPath := filepath.Join(filesFolderPath, fmt.Sprintf("%02x", i))
		if err := os.MkdirAll(subDirPath, pkg.DefaultDirPerm); err != nil {
			return fmt.Errorf("error creating files folders in path \"%s\": %s", subDirPath, err)
		}
	}

	// Create settings file
	if err := settings.Write(filepath.Join(path, settings.FileName), hashAlgorithm); err != nil {
		return fmt.Errorf("error creating settings file: %s", err)
	}
	return nil
}
