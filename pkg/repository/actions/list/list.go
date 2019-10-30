package list

import (
	"encoding/json"
	"fmt"
	api "github.com/Miguel-Dorta/gkup-backend/api/list"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
)

// snapshotNameRegex is the regex that represents the filename (and its parts) of a snapshot file.
var snapshotNameRegex = regexp.MustCompile("^(\\d{4})-(\\d{2})-(\\d{2})_(\\d{2})-(\\d{2})-(\\d{2}).json$")

// List takes a repository path and prints an api/list.SnapshotList object encoded in JSON to the outWriter provided.
// Errors will be printed to errWriter as strings separated by line-termination characters.
func List(path string, outWriter, errWriter io.Writer) {
	snapshots, err := getSnapshots(filepath.Join(path, repository.SnapshotsFolderName), errWriter)
	if err != nil {
		_, _ = fmt.Fprintln(errWriter, err)
		return
	}

	data, _ := json.Marshal(api.SnapshotList{SList: snapshots})
	_, _ = outWriter.Write(data)
}

// getSnapshots iterates a directory and returns a list of api/list.Snapshot objects.
// These objects will have the name of each subdirectory and will contain the times of their snapshots.
// An additional object with name=="" will be created for the times of the snapshots in path.
func getSnapshots(path string, errWriter io.Writer) ([]*api.Snapshot, error) {
	files, err := utils.ListDir(path)
	if err != nil {
		return nil, fmt.Errorf("error listing \"%s\": %w", path, err)
	}
	snapshots := make([]*api.Snapshot, 0, len(files))

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		times, err := getTimes(filepath.Join(path, file.Name()))
		if err != nil {
			_, _ = fmt.Fprintln(errWriter, err)
		}

		snapshots = append(snapshots, &api.Snapshot{
			Name:  file.Name(),
			Times: times,
		})
	}

	noNameTimes, err := getTimes(path)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return append(snapshots, &api.Snapshot{
		Name:  "",
		Times: noNameTimes,
	}), nil
}

// getTimes iterates the provided path and returns a list of api/list.Times objects
func getTimes(path string) ([]*api.Time, error) {
	files, err := utils.ListDir(path)
	if err != nil {
		return nil, fmt.Errorf("error listing \"%s\": %w", path, err)
	}

	times := make([]*api.Time, 0, len(files))

	for _, file := range files {
		if !file.Mode().IsRegular() || !snapshotNameRegex.MatchString(file.Name()) {
			continue
		}

		dateStrs := snapshotNameRegex.FindStringSubmatch(file.Name())[1:]
		Y, _ := strconv.Atoi(dateStrs[0])
		M, _ := strconv.Atoi(dateStrs[1])
		D, _ := strconv.Atoi(dateStrs[2])
		h, _ := strconv.Atoi(dateStrs[3])
		m, _ := strconv.Atoi(dateStrs[4])
		s, _ := strconv.Atoi(dateStrs[5])

		times = append(times, &api.Time{
			Year:   Y,
			Month:  M,
			Day:    D,
			Hour:   h,
			Minute: m,
			Second: s,
		})
	}
	return times, nil
}
