package main

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository"
	"os"
)

const (
	keyAction          = "action"
	keyGroup           = "group"
	keyPaths           = "paths"
	keyRepoPath        = "repository-path"
	keyRestorationPath = "restore-to"
	keyTime            = "time"
)

func main() {
	args, err := getArgs()
	if err != nil {
		printError(err.Error())
	}

	repoPath := getString(args, keyRepoPath)
	if repoPath == "" {
		printError(keyRepoPath + " is empty")
	}

	switch getString(args, keyAction) {
	case "backup":
		repository.Backup(
			repoPath,
			getString(args, keyGroup),
			getStringSlice(args, keyPaths),
		)
	case "check":
		repository.Check(repoPath)
	case "create":
		repository.Create(repoPath)
	case "list":
		repository.List(repoPath)
	case "restore":
		repository.Restore(
			repoPath,
			getString(args, keyRestorationPath),
			getString(args, keyGroup),
			getInt64(args, keyTime),
		)
	}
}

func printError(s string) {
	_, _ = fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}
