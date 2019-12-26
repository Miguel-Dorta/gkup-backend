package custom

import (
	"fmt"
	"path/filepath"
	"time"
)

type metadata struct {
	Version string `json:"version"`
}

type file struct {
	RelativePath string `json:"relative-path"`
	Hash    string `json:"hash"`
	Size    int64  `json:"size"`
}

const (
	Type         = "custom"
	snapshotsDir = "snapshots"
)

// TODO make tests

func getPath(repoPath, groupName string, timestamp int64) string {
	path := filepath.Join(repoPath, snapshotsDir)
	if groupName != "" {
		path = filepath.Join(path, groupName)
	}

	t := time.Unix(timestamp, 0).UTC()
	return filepath.Join(path, fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d.gkup",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}
