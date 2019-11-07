package snapshots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Snapshot struct {
	Version string    `json:"version"`
	Files   Directory `json:"files"`
}

type Directory struct {
	Name  string       `json:"name"`
	Dirs  []*Directory `json:"directories"`
	Files []*File      `json:"files"`
}

type File struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

func Read(path string) (*Snapshot, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading snapshot \"%s\": %w", path, err)
	}

	s := &Snapshot{}
	if err = json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("error parsing snapshot \"%s\": %w", path, err)
	}
	return s, nil
}
