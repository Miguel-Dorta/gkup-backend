package custom

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/NewSnapshot"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/x"
	"os"
)

type Decoder struct {
	hasReadVersion bool
	f              *os.File
	d              *json.Decoder
}

func NewDecoder(repoPath, snapshotName, snapshotTime string) (*Decoder, error) {
	path := formatPath(repoPath, snapshotName, snapshotTime)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading snapshot file \"%s\": %w", path, err)
	}
	return &Decoder{
		f: f,
		d: json.NewDecoder(f),
	}, nil
}

func (d *Decoder) Decode(f *x.File) error {
	if !d.hasReadVersion {
		if err := d.d.Decode(&versionJSON{}); err != nil {
			return err
		}
		d.hasReadVersion = true
	}
	return d.d.Decode(f)
}

func (d *Decoder) More() bool {
	return d.d.More()
}

func (d *Decoder) Close() error {
	return d.f.Close()
}
