package backup

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/output"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	pathPkg "path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func Backup(repoPath, groupName string, paths []string, outWriter, errWriter io.Writer) {
	startTime := time.Now()
	status := output.NewStatus(4, outWriter)
	defer status.Stop()

	sett, err := settings.Read(filepath.Join(repoPath, settings.FileName))
	if err != nil {
		printError(errWriter, "error reading repository settings: %s", err)
		return
	}

	status.NewStep("Listing files", 0)
	fileList := listFiles(paths, errWriter)
	if len(fileList) == 0 {
		printError(errWriter, "omitting empty snapshot")
		return
	}

	status.NewStep("Hashing files", len(fileList))
	if err := hashFileList(fileList, sett, status); err != nil {
		printError(errWriter, "error hashing files: %s", err)
		return
	}

	status.NewStep("Adding files to repository", len(fileList))
	if err := addFiles(fileList, repoPath, sett, status); err != nil {
		printError(errWriter, "error adding files to repository: %s", err)
		return
	}

	status.NewStep("Writing snapshot", 0)
	if err := writeSnapshot(fileList, repoPath, groupName, sett, startTime); err != nil {
		printError(errWriter, "error writing snapshot (it may be incomplete or corrupt): %s", err)
		return
	}
}

func listFiles(paths []string, errWriter io.Writer) []*file {
	list := make([]*file, 0, len(paths))

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			printError(errWriter, "error getting info from path \"%s\": %s", path, err)
			continue
		}

		if stat.IsDir() {
			list = append(list, listFilesRecursive(path, stat.Name(), errWriter)...)
		} else if stat.Mode().IsRegular() {
			list = append(list, &file{
				RealPath: path,
				File: files.File{
					RelativePath: stat.Name(),
					Hash:         nil,
					Size:         stat.Size(),
				},
			})
		} else {
			printError(errWriter, "file type of \"%s\" is not supported", path) //TODO
		}
	}

	return list
}

func listFilesRecursive(pathReal, pathRelative string, errWriter io.Writer) []*file {
	list, err := utils.ListDir(pathReal)
	if err != nil {
		printError(errWriter, "error listing directory \"%s\": %s", pathReal, err)
		return nil
	}

	fileList := make([]*file, 0, len(list))
	for _, f := range list {
		fPathReal := filepath.Join(pathReal, f.Name())
		fPathRelative := pathPkg.Join(pathRelative, f.Name())

		if f.IsDir() {
			fileList = append(fileList, listFilesRecursive(fPathReal, fPathRelative, errWriter)...)
		} else if f.Mode().IsRegular() {
			fileList = append(fileList, &file{
				RealPath: fPathReal,
				File: files.File{
					RelativePath: fPathRelative,
					Hash:         nil,
					Size:         f.Size(),
				},
			})
		} else {
			printError(errWriter, "file type of \"%s\" is not supported", fPathReal) //TODO
		}
	}

	return fileList
}

func hashFileList(l []*file, s *settings.Settings, status *output.Status) error {
	safeL := &safeFileList{list: l}
	err := new(error)
	wg := new(sync.WaitGroup)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			fileHasher(safeL, s, err, status)
			wg.Done()
		}()
	}
	wg.Wait()

	return *err
}

func fileHasher(l *safeFileList, s *settings.Settings, commonErr *error, status *output.Status) {
	hasher, _ := hash.NewHasher(s)
	for {
		f := l.next()
		if f == nil || *commonErr != nil {
			return
		}

		h, err := hasher.HashFile(f.RealPath)
		if err != nil {
			*commonErr = err
			return
		}

		f.Hash = h
		status.AddPart()
	}
}

func addFiles(l []*file, repoPath string, s *settings.Settings, status *output.Status) error {
	filesPath := filepath.Join(repoPath, files.FolderName)
	buf := make([]byte, s.BufferSize)

	for _, f := range l {
		if err := files.Add(filesPath, f.RealPath, &f.File, buf); err != nil {
			return err
		}
		status.AddPart()
	}
	return nil
}

func writeSnapshot(l []*file, repoPath, groupName string, s *settings.Settings, t time.Time) error {
	w, err := snapshots.NewWriter(repoPath, groupName, s, t)
	if err != nil {
		return fmt.Errorf("error creating snapshot writer: %s", err)
	}
	defer w.Close()

	for _, f := range l {
		if err := w.Write(&f.File); err != nil {
			return fmt.Errorf("error writing file (%s) to snapshot: %s", f.RelativePath, err)
		}
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing snapshot writer: %s", err)
	}
	return nil
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format+"\n", a...)
}
