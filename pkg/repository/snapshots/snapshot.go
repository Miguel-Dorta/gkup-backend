package snapshots

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/custom"
	"time"
)

type Reader interface {
	ReadNext() (*File, error)
	More() bool
	Close() error
}

type Writer interface {
	Write(f *File) error
	Close() error
}

type Snapshot struct {
	Group string
	Date int64
}

type File struct {
	RelPath string
	Hash []byte
	Size int64
}

func ListSnapshots(repoPath string, s *settings.Settings) ([]Snapshot, error) {
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
