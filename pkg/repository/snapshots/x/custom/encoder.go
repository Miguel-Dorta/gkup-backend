package custom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/internal"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/NewSnapshot"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/x"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"time"
)

type Encoder struct {
	hasWroteVersion bool
	f               *os.File
	e               *json.Encoder
}

func NewEncoder(repoPath, snapshotName string, t time.Time) (*Encoder, error) {
	path := formatPath(repoPath, snapshotName, formatFilename(t))

	f, err := utils.CreateWithParents(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create snapshot file \"%s\": %w", path, err)
	}
	return &Encoder{
		f: f,
		e: json.NewEncoder(f),
	}, nil
}

func (e *Encoder) Encode(f *x.File) error {
	if !e.hasWroteVersion {
		if err := e.e.Encode(&versionJSON{Version: internal.Version}); err != nil {
			return err
		}
		e.hasWroteVersion = true
	}
	return e.e.Encode(&fileJSON{
		RelPath: f.RelPath,
		Hash:    hex.EncodeToString(f.Hash),
		Size:    f.Size,
	})
}

func (e *Encoder) Close() error {
	return e.f.Close()
}
