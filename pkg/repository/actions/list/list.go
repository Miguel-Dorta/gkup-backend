package list

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"io"
	"path/filepath"
)

// TODO review tests

// List takes a repository path and prints an api/list.SnapshotList object encoded in JSON to the outWriter provided.
// Errors will be printed to errWriter as strings separated by line-termination characters.
func List(path string, outWriter, errWriter io.Writer) {
	sett, err := settings.Read(filepath.Join(path, settings.FileName))
	if err != nil {
		printError(errWriter, "error reading repository settings: %s", err)
		return
	}

	list, err := snapshots.List(path, sett)
	if err != nil {
		printError(errWriter, "error listing snapshots: %s", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		printError(errWriter, "error serializing snapshots: %s", err)
		return
	}
	_, _ = outWriter.Write(data)
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format + "\n", a...)
}
