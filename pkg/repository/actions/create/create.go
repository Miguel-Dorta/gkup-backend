package create

import (
	"encoding/hex"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO check tests

func Create(path string, s *settings.Settings, errWriter io.Writer) {
	// Check repo path
	if err := checkRepoPath(path); err != nil {
		printError(errWriter, err.Error())
		return
	}

	// Create files directories
	if err := createFilesDirs(path); err != nil {
		printError(errWriter, err.Error())
		return
	}

	// Create settings
	updateWithValidSettings(s)
	if err := settings.Write(filepath.Join(path, settings.FileName), s); err != nil {
		printError(errWriter, "error creating settings file: %s", err)
		return
	}
}

func checkRepoPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error getting repository path (%s) info: %s", path, err)
		}

		// If the error is that path doesn't exists, create it
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error creating repository directory (%s): %s", path, err)
		}
		return nil
	}

	if !stat.IsDir() {
		return fmt.Errorf("repository path (%s) is not a directory", path)
	}

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error listing repository path (%s) content: %s", path, err)
	}
	if len(fs) != 0 {
		return fmt.Errorf("repository path (%s) must be empty", path)
	}
	return nil
}

func createFilesDirs(path string) error {
	path = filepath.Join(path, files.FolderName)
	for i:=0x0; i<=0xff; i++ {
		iPath := filepath.Join(path, hex.EncodeToString([]byte{byte(i)}))
		if err := os.MkdirAll(iPath, 0755); err != nil {
			return fmt.Errorf("error creating directory \"%s\": %s", iPath, err)
		}
	}
	return nil
}

func printError(w io.Writer, format string, a ...interface{}) {
	if a == nil {
		_, _ = fmt.Fprintln(w, format)
		return
	}
	_, _ = fmt.Fprintf(w, format+"\n", a...)
}
