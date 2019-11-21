package custom

import (
	"fmt"
	"path/filepath"
	"time"
)

type versionJSON struct {
	Version string `json:"version"`
}

type fileJSON struct {
	RelPath string `json:"relative-path"`
	Hash    string `json:"hash"`
	Size    int64  `json:"size"`
}

const (
	snapshotsPath = "snapshots"
)

func formatFilename(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d.gkup",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func formatPath(repoPath, snapshotName, snapshotTime string) string {
	path := filepath.Join(repoPath, snapshotsPath)
	if snapshotName != "" {
		path = filepath.Join(path, snapshotName)
	}
	return filepath.Join(path, snapshotTime)
}
