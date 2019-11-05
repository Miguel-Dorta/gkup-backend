package restore

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type SnapshotToRestore struct {
	Name string
	Time struct {
		Year   int
		Month  int
		Day    int
		Hour   int
		Minute int
		Second int
	}
}

type copy struct {
	from, to string
}

func Restore(repoPath, restorePath string, bufferSize int, restoreSnap *SnapshotToRestore, outWriter, errWriter io.Writer) {
	destinyFiles, err := utils.ListDir(restorePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			_, _ = fmt.Fprintf(errWriter, "error listing destiny path: %s\n", err)
			return
		}

		if err := os.MkdirAll(restorePath, 0755); err != nil {
			_, _ = fmt.Fprintf(errWriter, "error creating destiny path: %s\n", err)
			return
		}
		destinyFiles = nil
	}

	if len(destinyFiles) != 0 {
		_, _ = fmt.Fprintln(errWriter, "destiny path is not empty")
		return
	}

	snap, err := getSnapshot(filepath.Join(repoPath, snapshots.FolderName), restoreSnap)
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "error getting snapshot: %s", err)
		return
	}

	copyList := getCopyList(repoPath, restorePath, snap)

	var progress int
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		statusPrinter(len(copyList), &progress, outWriter)
		wg.Done()
	}()

	buf := make([]byte, bufferSize)
	for progress = 0; progress < len(copyList); progress++ {
		c := copyList[progress]
		if err := utils.CopyFile(c.from, c.to, buf); err != nil {
			_, _ = fmt.Fprintf(errWriter, "error restoring file (%s): %s\n", c.to, err)
			continue
		}
	}

	wg.Wait()
}

func getSnapshot(path string, restoreSnap *SnapshotToRestore) (*snapshots.Snapshot, error) {
	if restoreSnap.Name != "" {
		path = filepath.Join(path, restoreSnap.Name)
	}
	path = filepath.Join(path, fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d.json",
		restoreSnap.Time.Year, restoreSnap.Time.Month, restoreSnap.Time.Day,
		restoreSnap.Time.Hour, restoreSnap.Time.Minute, restoreSnap.Time.Second))

	return snapshots.Read(path)
}

func getCopyList(repoPath, restorePath string, s *snapshots.Snapshot) []copy {
	return getCopyListRecursive(filepath.Join(repoPath, files.FolderName), restorePath, &s.Files)
}

func getCopyListRecursive(filesPath, relRestorePath string, d *snapshots.Directory) []copy {
	list := make([]copy, 0, len(d.Files))
	for _, f := range d.Files {
		list = append(list, copy{
			from: filepath.Join(filesPath, f.Hash[:2], f.Hash + "-" + strconv.FormatInt(f.Size, 10)),
			to:   filepath.Join(relRestorePath, f.Name),
		})
	}
	for _, subD := range d.Dirs {
		list = append(list, getCopyListRecursive(filesPath, filepath.Join(relRestorePath, d.Name), subD)...)
	}
	return list
}
