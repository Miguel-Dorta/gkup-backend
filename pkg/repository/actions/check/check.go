package check

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/output"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"io"
	"path/filepath"
	"runtime"
	"sync"
)

// TODO redo tests

func Check(repoPath string, outWriter, errWriter io.Writer) {
	status := output.NewStatus(0, outWriter)
	defer status.Stop()

	sett, err := settings.Read(filepath.Join(repoPath, settings.FileName))
	if err != nil {
		printError(errWriter, "error reading repository settings: %s", err)
		return
	}

	status.NewStep("Listing files", 0)
	list, err := files.List(filepath.Join(repoPath, files.FolderName))
	if err != nil {
		printError(errWriter, "error listing files: %s", err)
		return
	}

	status.NewStep("Checking files", len(list))
	checkFiles(list, sett, errWriter, status)
}

func checkFiles(list []string, s *settings.Settings, errWriter io.Writer, status *output.Status) {
	safeList := threadSafe.NewStringList(list)

	wg := new(sync.WaitGroup)
	for i:=0; i<runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			h, _ := hash.NewHasher(s)
			for {
				path := safeList.Next()
				if path == nil {
					return
				}

				if err := files.Check(*path, h); err != nil {
					printError(errWriter, "error checking file \"%s\": %s", *path, err)
				}
				status.AddPart()
			}
		}()
	}
	wg.Wait()
}

func printError(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format+"\n", a...)
}
