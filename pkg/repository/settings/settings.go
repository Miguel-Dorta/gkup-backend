package settings

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

//TODO add new stuff to tests and comment

const FileName = "settings.toml"

type Settings struct {
	Version       string `toml:"version"`
	BufferSize    int    `toml:"buffer_size"`
	HashAlgorithm string `toml:"hash_algorithm"`
	SnapshotType  string `toml:"snapshot_type"`
	DB            DB     `toml:"database"`
}

type DB struct {
	Host   string `toml:"hostname"`
	DBName string `toml:"database_name"`
	User   string `toml:"username"`
	Pass   string `toml:"password"`
	Port   int    `toml:"port"`
}

func Read(path string) (*Settings, error) {
	// Read data
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading setting in path \"%s\": %s", path, err)
	}

	// Parse settings
	s := new(Settings)
	if err := toml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error parsing settings: %s", err)
	}

	// Check if the fields have info
	if s.Version == "" || s.HashAlgorithm == "" || s.SnapshotType == "" {
		return nil, errors.New("incomplete information in settings")
	}
	return s, nil
}

func Write(path string, s *Settings) error {
	// Serialize settings
	data, err := toml.Marshal(s)
	if err != nil {
		return fmt.Errorf("error serializing settings: %s", err)
	}

	// Write data
	if err := ioutil.WriteFile(path, data, pkg.DefaultFilePerm); err != nil {
		return fmt.Errorf("error writing settings in path \"%s\": %s", path, err)
	}
	return nil
}
