package create

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/internal"
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	"path/filepath"
)

// TODO check tests

func Create(path string, errWriter io.Writer) {
	// Check existence
	if exists, err := fileutils.FileExist(path); err != nil {
		printError(errWriter, "error checking existence of path \"%s\": %s", path, err)
		return
	} else if !exists {
		if err := os.MkdirAll(path, pkg.DefaultDirPerm); err != nil {
			printError(errWriter, "error creating repo directory in path \"%s\": %s", path, err)
			return
		}
	}

	// Check if it's not a dir
	if stat, err := os.Stat(path); err != nil {
		printError(errWriter, "cannot get information from path \"%s\": %s", path, err)
		return
	} else if !stat.IsDir() {
		printError(errWriter, "repository path must be a directory")
		return
	}

	// Check if dir is empty
	if list, err := fileutils.ListDir(path); err != nil {
		printError(errWriter, "error listing repository directory (%s): %s\n", path, err)
		return
	} else if len(list) != 0 {
		printError(errWriter, "repository path must be empty")
		return
	}

	// Create settings
	if err := settings.Write(filepath.Join(path, settings.FileName), &settings.Settings{
		Version:       internal.Version,
		BufferSize:    128 * 1024,
		HashAlgorithm: "sha256",
		SnapshotType:  "custom",
		DB: settings.DB{
			Host:   "localhost",
			DBName: "gkup",
			User:   "user",
			Pass:   "pass",
			Port:   3306,
		},
	}); err != nil {
		printError(errWriter, "error creating settings file: %s", err)
		return
	}
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format + "\n", a...)
}
