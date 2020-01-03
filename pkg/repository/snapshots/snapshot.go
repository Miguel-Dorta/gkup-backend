package snapshots

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/custom"
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

func List(repoPath string, s *settings.Settings) (map[string][]int64, error) {
	switch s.SnapshotType {
	case custom.Type:
		return custom.List(repoPath, s)
	}
	return nil, fmt.Errorf("invalid snapshot type: %s", s.SnapshotType)
}

func NewReader(repoPath, groupName string, t int64, s *settings.Settings) (Reader, error) {
	switch s.SnapshotType {
	case custom.Type:
		return custom.NewReader(repoPath, groupName, t, s)
	}
	return nil, fmt.Errorf("invalid snapshot type: %s", s.SnapshotType)
}

func NewWriter(repoPath, groupName string, t int64, s *settings.Settings) (Writer, error) {
	switch s.SnapshotType {
	case custom.Type:
		return custom.NewWriter(repoPath, groupName, t, s)
	}
	return nil, fmt.Errorf("invalid snapshot type: %s", s.SnapshotType)
}
