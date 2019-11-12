package backup

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/internal"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"path/filepath"
	"time"
)

func getSnapshotPath(repoPath, snapName string, t time.Time) string {
	snapPath := filepath.Join(repoPath, snapshots.FolderName)
	if snapName != "" {
		snapPath = filepath.Join(snapPath, snapName)
	}
	snapPath = filepath.Join(snapPath, fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d.json",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
	return snapPath
}

func getSnapshot(pathsToBackup []string) (*snapshots.Snapshot, []*copyInfo, error) {
	s := &snapshots.Snapshot{
		Version: internal.Version,
		Files: snapshots.Directory{
			Name:  ".",
			Dirs:  make([]*snapshots.Directory, 0, len(pathsToBackup)),
			Files: make([]*snapshots.File, 0, len(pathsToBackup)),
		},
	}
	ciList := make([]*copyInfo, 0, 1000)

	for _, path := range pathsToBackup {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting information from path \"%s\": %w", path, err)
		}

		if stat.IsDir() {
			d, subList, err := getDir(path)
			if err != nil {
				return nil, nil, err
			}
			s.Files.Dirs = append(s.Files.Dirs, d)
			ciList = append(ciList, subList...)

		} else if stat.Mode().IsRegular() {
			f, err := getFile(path)
			if err != nil {
				return nil, nil, err
			}
			s.Files.Files = append(s.Files.Files, f.f)
			ciList = append(ciList, f)

		} else {
			return nil, nil, fmt.Errorf("unsupported file (%s)", path)
		}
	}

	return s, ciList, nil
}

func getDir(path string) (*snapshots.Directory, []*copyInfo, error) {
	fList, err := utils.ListDir(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot list directory \"%s\": %w", path, err)
	}

	ciList := make([]*copyInfo, 0, len(fList))
	d := &snapshots.Directory{
		Name:  filepath.Base(path),
		Dirs:  make([]*snapshots.Directory, 0, len(fList)),
		Files: make([]*snapshots.File, 0, len(fList)),
	}
	for _, f := range fList {
		fPath := filepath.Join(path, f.Name())

		if f.IsDir() {
			subD, subList, err := getDir(fPath)
			if err != nil {
				return nil, nil, err
			}
			d.Dirs = append(d.Dirs, subD)
			ciList = append(ciList, subList...)

		} else if f.Mode().IsRegular() {
			subF, err := getFile(fPath)
			if err != nil {
				return nil, nil, err
			}
			d.Files = append(d.Files, subF.f)
			ciList = append(ciList, subF)

		} else {
			return nil, nil, fmt.Errorf("unsupported file (%s)", fPath)
		}
	}
	return d, ciList, nil
}

func getFile(path string) (*copyInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get information from file \"%s\": %w", path, err)
	}

	return &copyInfo{
		f: &snapshots.File{
			Name: stat.Name(),
			Hash: "",
			Size: stat.Size(),
		},
		path: path,
	}, nil
}
