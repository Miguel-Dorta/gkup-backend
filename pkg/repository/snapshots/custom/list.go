package custom

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func List(repoPath string, _ *settings.Settings) ([]snapshots.Snapshot, error) {
	snapshotsPath := filepath.Join(repoPath, snapshotsDir)
	list, err := utils.ListDir(snapshotsPath)
	if err != nil {
		return nil, fmt.Errorf("cannot list snapshot dir (%s): %s", snapshotsPath, err)
	}
	return getSnapshots(filepath.Join(repoPath, snapshotsDir), list), nil
}

func getSnapshots(snapshotsPath string, list []os.FileInfo) []snapshots.Snapshot {
	snapList := make([]snapshots.Snapshot, 0, len(list))
	noGroupSnaps := make([]os.FileInfo, 0, len(list))

	for _, f := range list {
		if f.IsDir() {
			subList, err := utils.ListDir(filepath.Join(snapshotsPath, f.Name()))
			if err != nil {
				continue
			}
			snapList = append(snapList, getGroupSnapshots(f.Name(), subList)...)
		} else if f.Mode().IsRegular() {
			noGroupSnaps = append(noGroupSnaps, f)
		}
	}

	snapList = append(snapList, getGroupSnapshots("", noGroupSnaps)...)
	sort.Slice(snapList, func(i, j int) bool {
		if snapList[i].Group != snapList[j].Group {
			return snapList[i].Group < snapList[j].Group
		}
		return snapList[i].Date < snapList[j].Date
	})
	return snapList
}

func getGroupSnapshots(groupName string, list []os.FileInfo) []snapshots.Snapshot {
	snapList := make([]snapshots.Snapshot, 0, len(list))

	for _, f := range list {
		if !f.Mode().IsRegular() {
			continue
		}

		t := getDateFromFilename(f.Name())
		if t == nil {
			continue
		}

		snapList = append(snapList, snapshots.Snapshot{
			Group: groupName,
			Date:  t.Unix(),
		})
	}

	return snapList
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
