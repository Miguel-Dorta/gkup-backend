package snapshots

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/custom"
	"time"
)

type Reader interface {
	ReadNext() (*files.File, error)
	More() bool
	Close() error
}

type Writer interface {
	Write(f *files.File) error
	Close() error
}

func IsValidType(s string) bool {
	switch s {
	case custom.Type:
		return true
	}
	return false
}

func List(repoPath string, s *settings.Settings) (map[string][]int64, error) {
	switch s.SnapshotType {
	case custom.Type:
		return custom.List(repoPath, s)
	}
	return nil, fmt.Errorf("invalid snapshot type: %s", s.SnapshotType)
}

func NewReader(repoPath, groupName string, s *settings.Settings, t time.Time) (Reader, error) {
	switch s.SnapshotType {
	case custom.Type:
		return custom.NewReader(repoPath, groupName, s, t)
	}
	return nil, fmt.Errorf("invalid snapshot type: %s", s.SnapshotType)
}

func NewWriter(repoPath, groupName string, s *settings.Settings, t time.Time) (Writer, error) {
	switch s.SnapshotType {
	case custom.Type:
		return custom.NewWriter(repoPath, groupName, s, t)
	}
	return nil, fmt.Errorf("invalid snapshot type: %s", s.SnapshotType)
}
