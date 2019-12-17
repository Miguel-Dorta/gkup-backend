package check

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	total    int
	progress = new(threadSafe.Counter)
)

func Check(path string, outWriter, errWriter io.Writer) {
	// Get hash algorithm
	sett, err := settings.Read(filepath.Join(path, settings.FileName))
	if err != nil {
		printError(errWriter, "error reading repository settings: %s", err)
		return
	}
	if !hash.IsValidHashAlgorithm(sett.HashAlgorithm) {
		printError(errWriter, "invalid hash algorithm found in settings: %s", sett.HashAlgorithm)
		return
	}

	// Get file list
	fileList := listFiles(filepath.Join(path, files.FolderName), errWriter)
	if fileList == nil || len(fileList) == 0 {
		return
	}

	// Start status printer routine
	total = len(fileList)
	statusFinished := make(chan bool, 1)
	go func() {
		statusPrinter(outWriter)
		statusFinished <- true
	}()

	// Start file checkers
	safeFileList := threadSafe.NewStringList(fileList)
	wg := new(sync.WaitGroup)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			fileChecker(safeFileList, sett, errWriter)
			wg.Done()
		}()
	}

	// Wait and close the execution
	wg.Wait()
	<- statusFinished
}

func fileChecker(list *threadSafe.StringList, s *settings.Settings, errWriter io.Writer) {
	h, _ := hash.NewHasher(s)

	for {
		path := list.Next()
		if path == nil {
			return
		}

		if err := files.Check(*path, h); err != nil {
			printError(errWriter, "error checking file \"%s\": %s", *path, err)
		}
		progress.Add(1)
	}
}

func listFiles(path string, errWriter io.Writer) []string {
	fileList := make([]string, 0, 1000)

	path = filepath.Join(path, files.FolderName)
	dirs, err := utils.ListDir(path)
	if err != nil {
		printError(errWriter, "cannot list files directory: %s", err)
		return nil
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dirPath := filepath.Join(path, dir.Name())

		fs, err := utils.ListDir(dirPath)
		if err != nil {
			printError(errWriter, "cannot list directory \"%s\": %s", dirPath, err)
			continue
		}

		for _, f := range fs {
			if !f.Mode().IsRegular() {
				continue
			}
			fileList = append(fileList, filepath.Join(dirPath, f.Name()))
		}
	}
	return fileList
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format+"\n", a...)
}
