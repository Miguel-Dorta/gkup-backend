package repository

import (
	"github.com/Miguel-Dorta/gkup-backend/api"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/backup"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/check"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/create"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/list"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/restore"
	"io"
)

func Backup(repoPath, snapName string, pathsToBackup []string, bufferSize, threads int, outWriter, errWriter io.Writer) {
	backup.Backup(repoPath, snapName, pathsToBackup, bufferSize, threads, outWriter, errWriter)
}

func Check(path string, threads, bufferSize int, outWriter, errWriter io.Writer) {
	check.Check(path, threads, bufferSize, outWriter, errWriter)
}

func Create(path, hashAlgorithm string, errWriter io.Writer) {
	create.Create(path, hashAlgorithm, errWriter)
}

func List(path string, outWriter, errWriter io.Writer) {
	list.List(path, outWriter, errWriter)
}

func Restore(repoPath, restorePath string, bufferSize int, restoreSnap *api.Snapshot, outWriter, errWriter io.Writer) {
	restore.Restore(repoPath, restorePath, bufferSize, restoreSnap, outWriter, errWriter)
}
