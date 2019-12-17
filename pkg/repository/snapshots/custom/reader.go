package custom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"io"
	"os"
	"time"
)

type Reader struct {
	f io.ReadCloser
	d *json.Decoder
}

func NewReader(repoPath, groupName string, s *settings.Settings, t time.Time) (*Reader, error) {
	path := getPath(repoPath, groupName, t)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading snapshot file \"%s\": %w", path, err)
	}
	bufF := newFileReader(f, s.BufferSize)

	r := &Reader{
		f: bufF,
		d: json.NewDecoder(bufF),
	}
	if err = r.d.Decode(&metadata{}); err != nil {
		return nil, fmt.Errorf("cannot read metadata from snapshot file (%s): %w", path, err)
	}
	return r, nil
}

func (r *Reader) ReadNext() (*snapshots.File, error) {
	var j file
	if err := r.d.Decode(j); err != nil {
		return nil, fmt.Errorf("error decoding file: %w", err)
	}

	hash, err := hex.DecodeString(j.Hash)
	if err != nil {
		return nil, fmt.Errorf("error decoding hash: %w", err)
	}

	return &snapshots.File{
		RelPath: j.RelPath,
		Hash:    hash,
		Size:    j.Size,
	}, nil
}

func (r *Reader) More() bool {
	return r.d.More()
}

func (r *Reader) Close() error {
	return r.f.Close()
}
