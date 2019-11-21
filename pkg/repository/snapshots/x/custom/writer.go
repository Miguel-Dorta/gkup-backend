package custom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/internal"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/x"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"time"
)

type Writer struct {
	f *os.File
	e *json.Encoder
}

func NewWriter(repoPath, snapshotName string, t time.Time) (*Writer, error) {
	path := formatPath(repoPath, snapshotName, formatFilename(t))

	f, err := utils.CreateWithParents(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create snapshot file \"%s\": %w", path, err)
	}
	w := &Writer{
		f: f,
		e: json.NewEncoder(f),
	}
	if err = w.e.Encode(&versionJSON{Version: internal.Version}); err != nil {
		return nil, fmt.Errorf("cannot write version to snapshot file (%s): %w", path, err)
	}
	return w, nil
}

func (w *Writer) Write(f *x.File) error {
	return w.e.Encode(&fileJSON{
		RelPath: f.RelPath,
		Hash:    hex.EncodeToString(f.Hash),
		Size:    f.Size,
	})
}

func (w *Writer) Close() error {
	return w.f.Close()
}
