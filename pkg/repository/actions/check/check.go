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
	"sync"
)

func Check(path string, threads, bufferSize int, outWriter, errWriter io.Writer) {
	// Get hash algorithm
	sett, err := settings.Read(filepath.Join(path, settings.FileName))
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "error reading repository settings: %s\n", err)
		return
	}
	if _, exists := hash.Algorithms[sett.HashAlgorithm]; !exists {
		_, _ = fmt.Fprintf(errWriter, "invalid hash algorithm (%s)\n", sett.HashAlgorithm)
		return
	}

	// Get file list
	fileList := listFiles(filepath.Join(path, files.FolderName), errWriter)
	if fileList == nil || len(fileList) == 0 {
		return
	}

	// Start status printer routine
	progress := &threadSafe.Counter{}
	quit := make(chan bool)
	wgStatus := &sync.WaitGroup{}
	wgStatus.Add(1)
	go func() {
		statusPrinter(len(fileList), progress, outWriter, quit)
		wgStatus.Done()
	}()

	// Start file checkers
	safeFileList := threadSafe.NewStringList(fileList)
	wg := &sync.WaitGroup{}
	for i:=0; i<threads; i++ {
		wg.Add(1)
		go func() {
			fileChecker(safeFileList, progress, errWriter, sett.HashAlgorithm, bufferSize)
			wg.Done()
		}()
	}

	// Wait and close the execution
	wg.Wait()
	quit <- true
	wgStatus.Wait()
}

func fileChecker(list *threadSafe.StringList, progress *threadSafe.Counter, errWriter io.Writer, hashAlgorithm string, bufferSize int) {
	h := hash.Algorithms[hashAlgorithm]()
	buf := make([]byte, bufferSize)

	for {
		path := list.Next()
		if path == nil {
			return
		}

		if err := files.Check(*path, h, buf); err != nil {
			_, _ = fmt.Fprintf(errWriter,"error checking file \"%s\": %s\n", *path, err)
		}
		progress.Add(1)
	}
}

func listFiles(path string, errWriter io.Writer) []string {
	fileList := make([]string, 0, 256)

	dirs, err := utils.ListDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(errWriter, "cannot list files directory: %s\n", err)
		return nil
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dirPath := filepath.Join(path, dir.Name())

		fs, err := utils.ListDir(dirPath)
		if err != nil {
			_, _ = fmt.Fprintf(errWriter, "error listing files in directory \"%s\": %s\n", dirPath, err)
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
