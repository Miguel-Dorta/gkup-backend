package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository"
	"os"
)

type input struct {
	Action   string                 `json:"action"`
	RepoPath string                 `json:"repository-path"`
	Args     map[string]interface{} `json:"args"`
}

const (
	keyOutputTimeInMS = "output-time-ms"
	keyGroupName = "group"
	keyPaths = "paths"
	keyRestorationPath = "restore-to"
	keySnapTime = "snapshot-time"
)

func getInput() (input, error) {
	var in input
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return input{}, fmt.Errorf("error decoding json: %s", err)
	}

	if in.RepoPath == "" {
		return input{}, errors.New("undefined repository-path")
	}

	return in, nil
}

func main() {
	in, err := getInput()
	if err != nil {
		printError(err.Error())
		return // Just for the linter to stop complaining
	}

	switch in.Action {
	case "create":
		repository.Create(in.RepoPath)
	case "list":
		repository.List(in.RepoPath)
	case "backup":
		repository.Backup(
			in.RepoPath,
			getInt(in.Args, keyOutputTimeInMS),
			getString(in.Args, keyGroupName),
			getStringSlice(in.Args, keyPaths))
	case "check":
		repository.Check(
			in.RepoPath,
			getInt(in.Args, keyOutputTimeInMS))
	case "restore":
		repository.Restore(
			in.RepoPath,
			getInt(in.Args, keyOutputTimeInMS),
			getString(in.Args, keyRestorationPath),
			getString(in.Args, keyGroupName),
			getInt64(in.Args, keySnapTime))
	default:
		printError("invalid or undefined action")
	}
}

func printError(s string) {
	_, _ = fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}
