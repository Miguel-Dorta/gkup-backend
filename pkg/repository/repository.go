package repository

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/backup"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/check"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/create"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/list"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/restore"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"os"
)

var (
	stdout = threadSafe.NewWriter(os.Stdout)
	stderr = threadSafe.NewWriter(os.Stderr)
)

func Create(repoPath string) {
	create.Create(repoPath, stderr)
}

func List(repoPath string) {
	list.List(repoPath, stdout, stderr)
}

func Backup(repoPath string, outputTimeInMS int, groupName string, paths []string) {
	backup.Backup(repoPath, outputTimeInMS, groupName, paths, stdout, stderr)
}

func Check(repoPath string, outputTimeInMS int) {
	check.Check(repoPath, outputTimeInMS, stdout, stderr)
}

func Restore(repoPath string, outputTimeInMS int, restorationPath, snapGroup string, snapTime int64) {
	restore.Restore(repoPath, outputTimeInMS, restorationPath, snapGroup, snapTime, stdout, stderr)
}
