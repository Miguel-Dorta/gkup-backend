package custom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"io"
	"os"
)

type Reader struct {
	f io.ReadCloser
	d *json.Decoder
}

func NewReader(repoPath, groupName string, t int64, s *settings.Settings) (*Reader, error) {
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

func (r *Reader) ReadNext() (*files.File, error) {
	var j file
	if err := r.d.Decode(j); err != nil {
		return nil, fmt.Errorf("error decoding file: %w", err)
	}

	hash, err := hex.DecodeString(j.Hash)
	if err != nil {
		return nil, fmt.Errorf("error decoding hash: %w", err)
	}

	return &files.File{
		RelativePath: j.RelativePath,
		Hash:         hash,
		Size:         j.Size,
	}, nil
}

func (r *Reader) More() bool {
	return r.d.More()
}

func (r *Reader) Close() error {
	return r.f.Close()
}
