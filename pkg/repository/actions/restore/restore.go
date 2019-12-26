package restore

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/Miguel-Dorta/gkup-backend/pkg/output"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	"path/filepath"
)

// TODO rewrite tests

func Restore(repoPath, restorationPath, snapGroup string, snapTime int64, outWriter, errWriter io.Writer) {
	status := output.NewStatus(2, outWriter)
	defer status.Stop()

	sett, err := settings.Read(filepath.Join(repoPath, settings.FileName))
	if err != nil {
		printError(errWriter, "error reading repository settings: %s", err)
		return
	}

	if err := checkRestorationPath(restorationPath); err != nil {
		printError(errWriter, err.Error())
		return
	}

	status.NewStep("Listing files", 0)
	list, err := listFiles(repoPath, snapGroup, snapTime, sett)
	if err != nil {
		printError(errWriter, "error listing files: %s", err)
		return
	}

	status.NewStep("Restoring files", len(list))
	restoreFiles(list, repoPath, restorationPath, errWriter, status)
}

func listFiles(repoPath, snapGroup string, snapTime int64, s *settings.Settings) ([]*files.File, error) {
	list := make([]*files.File, 0, 1000)

	r, err := snapshots.NewReader(repoPath, snapGroup, snapTime, s)
	if err != nil {
		return nil, fmt.Errorf("error creating snapshot reader: %s", err)
	}
	defer r.Close()

	for r.More() {
		f, err := r.ReadNext()
		if err != nil {
			return nil, fmt.Errorf("error listing file: %s", err)
		}
		list = append(list, f)
	}

	return list, nil
}

func restoreFiles(l []*files.File, repoPath, restorationPath string, errWriter io.Writer, status *output.Status) {
	filesPath := filepath.Join(repoPath, files.FolderName)
	b := make([]byte, 128*1024)

	for _, f := range l {
		if err := files.Restore(filesPath, restorationPath, f, b); err != nil {
			printError(errWriter, "error restoring file \"%s\": %s", f.RelativePath, err)
		}
		status.AddPart()
	}
}

func checkRestorationPath(path string) error {
	list, err := utils.ListDir(path)
	if err == nil {
		if len(list) != 0 {
			return errors.New("restoration path is not empty")
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("error listing restoration path: %s", err)
	}

	if err = os.MkdirAll(path, pkg.DefaultDirPerm); err != nil {
		return fmt.Errorf("error creating restoration path: %s", err)
	}
	return nil
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format + "\n", a...)
}
