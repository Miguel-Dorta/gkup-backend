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
	RelPath string `json:"relative-path"`
	Hash    string `json:"hash"`
	Size    int64  `json:"size"`
}

const (
	Type         = "custom"
	bufferSize   = 128*1024
	snapshotsDir = "snapshots"
)

func getPath(repoPath, groupName string, t time.Time) string {
	path := filepath.Join(repoPath, snapshotsDir)
	if groupName != "" {
		path = filepath.Join(path, groupName)
	}
	return filepath.Join(path, fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d.gkup",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}