package custom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/internal"
	"github.com/Miguel-Dorta/gkup-backend/pkg/fileutils"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"io"
)

type Writer struct {
	f io.WriteCloser
	e *json.Encoder
}

func NewWriter(repoPath, groupName string, t int64, s *settings.Settings) (*Writer, error) {
	path := getPath(repoPath, groupName, t)

	f, err := fileutils.CreateWithParents(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create snapshot file \"%s\": %w", path, err)
	}
	bufF := newFileWriter(f, s.BufferSize)

	w := &Writer{
		f: bufF,
		e: json.NewEncoder(bufF),
	}
	if err = w.e.Encode(&metadata{Version: internal.Version}); err != nil {
		return nil, fmt.Errorf("cannot write metadata to snapshot file (%s): %w", path, err)
	}
	return w, nil
}

func (w *Writer) Write(f *files.File) error {
	return w.e.Encode(&file{
		RelativePath: f.RelativePath,
		Hash:         hex.EncodeToString(f.Hash),
		Size:         f.Size,
	})
}

func (w *Writer) Close() error {
	return w.f.Close()
}
