package custom

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/fileutils"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func List(repoPath string, _ *settings.Settings) (map[string][]int64, error) {
	snapshotsPath := filepath.Join(repoPath, snapshotsDir)
	list, err := fileutils.ListDir(snapshotsPath)
	if err != nil {
		return nil, fmt.Errorf("cannot list snapshot dir (%s): %s", snapshotsPath, err)
	}
	return getSnapshots(filepath.Join(repoPath, snapshotsDir), list), nil
}

func getSnapshots(snapshotsPath string, list []os.FileInfo) map[string][]int64 {
	snaps := make(map[string][]int64, len(list))

	for _, f := range list {
		if f.IsDir() {
			subList, err := fileutils.ListDir(filepath.Join(snapshotsPath, f.Name()))
			if err != nil {
				continue
			}
			snaps[f.Name()] = getTimes(subList)
		} else if f.Mode().IsRegular() {
			t := getDateFromFilename(f.Name())
			if t == nil {
				continue
			}
			snaps[""] = append(snaps[""], t.Unix())
		}
	}

	return snaps
}

func getTimes(list []os.FileInfo) []int64 {
	times := make([]int64, 0, len(list))

	for _, f := range list {
		if !f.Mode().IsRegular() {
			continue
		}

		t := getDateFromFilename(f.Name())
		if t == nil {
			continue
		}

		times = append(times, t.Unix())
	}
	return times
}

func getDateFromFilename(name string) *time.Time {
	if len(name) != 24 {
		return nil
	}

	Y, err := strconv.Atoi(name[:4])
	if err != nil {
		return nil
	}
	M, err := strconv.Atoi(name[5:7])
	if err != nil {
		return nil
	}
	D, err := strconv.Atoi(name[8:10])
	if err != nil {
		return nil
	}
	h, err := strconv.Atoi(name[11:13])
	if err != nil {
		return nil
	}
	m, err := strconv.Atoi(name[14:16])
	if err != nil {
		return nil
	}
	s, err := strconv.Atoi(name[17:19])
	if err != nil {
		return nil
	}

	t := time.Date(Y, time.Month(M), D, h, m, s, 0, time.UTC)
	return &t
}
