package create

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	"path/filepath"
)

func Create(path string, s *settings.Settings, errWriter io.Writer) {
	log := logger{errWriter}

	// Check existence
	if exists, err := utils.FileExist(path); err != nil {
		log.errorf("error checking existence of path \"%s\": %s", path, err)
		return
	} else if !exists {
		if err := os.MkdirAll(path, pkg.DefaultDirPerm); err != nil {
			log.errorf("error creating repo directory in path \"%s\": %s", path, err)
			return
		}
	}

	// Check if it's not a dir
	if stat, err := os.Stat(path); err != nil {
		log.errorf("cannot get information from path \"%s\": %s", path, err)
		return
	} else if !stat.IsDir() {
		log.errorf("repository path must be a directory")
		return
	}

	// Check if dir is empty
	if list, err := utils.ListDir(path); err != nil {
		log.errorf("error listing repository directory (%s): %s\n", path, err)
		return
	} else if len(list) != 0 {
		log.errorf("repository path must be empty")
		return
	}

	// Create settings
	if err := settings.Write(filepath.Join(path, settings.FileName), s); err != nil {
		log.errorf("error creating settings file: %s", err)
		return
	}
}
