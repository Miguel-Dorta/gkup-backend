package backup

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// thread must be >= 1, bufferSize should be >= 512
func Backup(repoPath, snapName string, pathsToBackup []string, bufferSize, threads int, outWriter, errWriter io.Writer) {
	startTime := time.Now()

	// Get snapshot and list of files
	snap, ciList, err := getSnapshot(pathsToBackup)
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "error getting snapshot: %s\n", err)
		return
	}

	// Load settings and check that it has a valid hash algorithm
	sett, err := settings.Read(filepath.Join(repoPath, settings.FileName))
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "error reading repository settings: %s\n", err)
		return
	}
	if _, exists := hash.Algorithms[sett.HashAlgorithm]; !exists {
		_, _ = fmt.Fprintln(errWriter, "invalid hash algorithm found in repository settings")
		return
	}

	// Start status printer
	progress := new(threadSafe.Counter)
	quitStatus := startStatus(progress, len(ciList) * 2, outWriter)

	if err = hashFilesParallel(ciList, sett.HashAlgorithm, bufferSize, threads, progress); err != nil {
		_, _ = fmt.Fprintf(errWriter, "error hashing files: %s\n", err)
		quitStatus()
		return
	}

	// Add files to repo
	buf := make([]byte, bufferSize)
	filesPath := filepath.Join(repoPath, files.FolderName)
	for _, ci := range ciList {
		if err = addFileToRepo(ci, filesPath, buf); err != nil {
			_, _ = fmt.Fprintf(errWriter, "error adding file to repo: %s\n", err)
			continue
		}
		progress.Add(1)
	}

	// Stop status printer
	quitStatus()

	// Write snapshot
	if err := snapshots.Write(getSnapshotPath(repoPath, snapName, startTime), snap); err != nil {
		_, _ = fmt.Fprintf(errWriter, "\n")
	}
}

func addFileToRepo(ci *copyInfo, filesPath string, buf []byte) error {
	// Generate path where it should be saved
	destinyPath := filepath.Join(filesPath, ci.f.Hash[:2], ci.f.Hash + "-" + strconv.FormatInt(ci.f.Size, 10))

	// If exists (no error given) skip
	_, err := os.Stat(destinyPath)
	if err == nil {
		return nil
	}

	// If there's an error that is not a "Not exist" error, return error
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error getting info from file \"%s\": %w", destinyPath, err)
	}

	// Copy it where it should be
	if err = utils.CopyFile(ci.path, destinyPath, buf); err != nil {
		return fmt.Errorf("error copying file \"%s\" to repository: %w", ci.path, err)
	}
	return nil
}
