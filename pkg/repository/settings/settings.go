package settings

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash/checkHash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots/custom"
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
		return nil, fmt.Errorf("error reading setting in path \"%s\": %w", path, err)
	}

	// Parse settings
	s := new(Settings)
	if err := toml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error parsing settings: %w", err)
	}

	// Check fields
	if err = check(s); err != nil {
		return nil, fmt.Errorf("error checking settings: %w", err)
	}

	return s, nil
}

func Write(path string, s *Settings) error {
	// Check settings
	if err := check(s); err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}

	// Serialize settings
	data, err := toml.Marshal(s)
	if err != nil {
		return fmt.Errorf("error serializing settings: %w", err)
	}

	// Write data
	if err := ioutil.WriteFile(path, data, pkg.DefaultFilePerm); err != nil {
		return fmt.Errorf("error writing settings in path \"%s\": %w", path, err)
	}
	return nil
}

func check(s *Settings) error {
	if s.Version == "" {
		return errors.New("invalid or nonexistent version")
	}

	if s.BufferSize < 512 {
		return errors.New("buffer_size is too small or nonexistent")
	}

	if !checkHash.ValidAlgorithm(s.HashAlgorithm) {
		return errors.New("invalid or nonexistent hash_algorithm")
	}

	if !snapshots.IsValidType(s.SnapshotType) {
		return errors.New("invalid or nonexistent snapshot_type")
	}
	if s.SnapshotType == custom.Type {
		return nil
	}

	if s.DB.Port <= 0 || s.DB.Port > 0xFFFF {
		return errors.New("invalid or nonexistent database port")
	}
	switch "" {
	case s.DB.Host:
		return errors.New("empty database hostname")
	case s.DB.DBName:
		return errors.New("empty database name")
	case s.DB.User:
		return errors.New("empty database username")
	case s.DB.Pass:
		return errors.New("empty database password")
	}

	return nil
}
