package create

import (
	"github.com/Miguel-Dorta/gkup-backend/internal"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
)

const (
	defaultBufferSize = 128*1024
	defaultHashAlgorithm = "sha256"
	defaultSnapshotType = "custom"
	defaultDBHost = "localhost"
	defaultDBPort = 3306
	defaultDBUser = "user"
	defaultDBPass = "pass"
	defaultDBName = "gkup"
)

func updateWithValidSettings(s *settings.Settings) {
	s.Version = internal.Version

	if s.BufferSize == 0 {
		s.BufferSize = defaultBufferSize
	}
	if s.HashAlgorithm == "" {
		s.HashAlgorithm = defaultHashAlgorithm
	}
	if s.SnapshotType == "" {
		s.SnapshotType = defaultSnapshotType
	}
	if s.DB.Host == "" {
		s.DB.Host = defaultDBHost
	}
	if s.DB.Port == 0 {
		s.DB.Port = defaultDBPort
	}
	if s.DB.User == "" {
		s.DB.User = defaultDBUser
	}
	if s.DB.Pass == "" {
		s.DB.Pass = defaultDBPass
	}
	if s.DB.DBName == "" {
		s.DB.DBName = defaultDBName
	}
}
