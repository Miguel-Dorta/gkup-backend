package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"os"
)

type input struct {
	Action   string                 `json:"action"`
	RepoPath string                 `json:"repository-path"`
	Args     map[string]interface{} `json:"args"`
}

const (
	// Common keys
	keyOutputTimeInMS = "output-time-ms"
	keyGroupName = "group"

	// Backup keys
	keyPaths = "paths"

	// Restore keys
	keyRestorationPath = "restore-to"
	keySnapTime = "snapshot-time"

	// Settings keys
	keyBufferSize = "buffer-size"
	keyHashAlgorithm = "hash-algorithm"
	keySnapshotType = "snapshot-type"
	keyDBHost = "database-host"
	keyDBPort = "database-port"
	keyDBName = "database-name"
	keyDBUser = "database-user"
	keyDBPass = "database-password"
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
	case "create":
		repository.Create(
			in.RepoPath,
			&settings.Settings{
				BufferSize:    getOptionalInt(in.Args, keyBufferSize),
				HashAlgorithm: getOptionalString(in.Args, keyHashAlgorithm),
				SnapshotType:  getOptionalString(in.Args, keySnapshotType),
				DB:            settings.DB{
					Host:   getOptionalString(in.Args, keyDBHost),
					DBName: getOptionalString(in.Args, keyDBName),
					User:   getOptionalString(in.Args, keyDBUser),
					Pass:   getOptionalString(in.Args, keyDBPass),
					Port:   getOptionalInt(in.Args, keyDBPort),
				},
			})
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
