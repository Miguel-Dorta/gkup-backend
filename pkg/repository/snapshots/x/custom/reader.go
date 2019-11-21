package custom

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/x"
	"os"
)

type Reader struct {
	f *os.File
	d *json.Decoder
}

func NewReader(repoPath, snapshotName, snapshotTime string) (*Reader, error) {
	path := formatPath(repoPath, snapshotName, snapshotTime)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading snapshot file \"%s\": %w", path, err)
	}
	r := &Reader{
		f: f,
		d: json.NewDecoder(f),
	}
	if err = r.d.Decode(&versionJSON{}); err != nil {
		return nil, fmt.Errorf("cannot read version from snapshot file (%s): %w", path, err)
	}
	return r, nil
}

func (r *Reader) ReadNext() (f *x.File, err error) {
	err = r.d.Decode(f)
	return f, err
}

func (r *Reader) More() bool {
	return r.d.More()
}

func (r *Reader) Close() error {
	return r.f.Close()
}
