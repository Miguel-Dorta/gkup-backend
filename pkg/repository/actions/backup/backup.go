package backup

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/output"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"io"
	"path/filepath"
	"time"
)

func Backup(repoPath, groupName string, paths []string, outWriter, errWriter io.Writer) {
	startTime := time.Now()
	status := output.NewStatus(4, outWriter)
	defer status.Stop()

	sett, err := settings.Read(filepath.Join(repoPath, settings.FileName))
	if err != nil {
		printError(errWriter, "error reading repository settings: %s", err)
		return
	}

	status.NewStep("Listing files", 0)
	fileList := listFiles(paths, errWriter)
	if len(fileList) == 0 {
		printError(errWriter, "omitting empty snapshot")
		return
	}

	status.NewStep("Hashing files", len(fileList))
	if err := hashFileList(fileList, sett, status); err != nil {
		printError(errWriter, "error hashing files: %s", err)
		return
	}

	status.NewStep("Adding files to repository", len(fileList))
	if err := addFiles(fileList, repoPath, sett, status); err != nil {
		printError(errWriter, "error adding files to repository: %s", err)
		return
	}

	status.NewStep("Writing snapshot", 0)
	if err := writeSnapshot(fileList, repoPath, groupName, sett, startTime); err != nil {
		printError(errWriter, "error writing snapshot (it may be incomplete or corrupt): %s", err)
		return
	}
}

func addFiles(l []*file, repoPath string, s *settings.Settings, status *output.Status) error {
	filesPath := filepath.Join(repoPath, files.FolderName)
	buf := make([]byte, s.BufferSize)

	for _, f := range l {
		if err := files.Add(filesPath, f.RealPath, &f.File, buf); err != nil {
			return err
		}
		status.AddPart()
	}
	return nil
}

func writeSnapshot(l []*file, repoPath, groupName string, s *settings.Settings, t time.Time) error {
	w, err := snapshots.NewWriter(repoPath, groupName, t.UTC().Unix(), s)
	if err != nil {
		return fmt.Errorf("error creating snapshot writer: %s", err)
	}
	defer w.Close()

	for _, f := range l {
		if err := w.Write(&f.File); err != nil {
			return fmt.Errorf("error writing file (%s) to snapshot: %s", f.RelativePath, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing snapshot writer: %s", err)
	}
	return nil
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format+"\n", a...)
}
