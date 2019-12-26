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

func Backup(repoPath, groupName string, paths []string) {
	backup.Backup(repoPath, groupName, paths, stdout, stderr)
}

func Check(repoPath string) {
	check.Check(repoPath, stdout, stderr)
}

func Create(repoPath string) {
	create.Create(repoPath, stderr)
}

func List(repoPath string) {
	list.List(repoPath, stdout, stderr)
}

func Restore(repoPath, restorationPath, snapGroup string, snapTime int64) {
	restore.Restore(repoPath, restorationPath, snapGroup, snapTime, stdout, stderr)
}
