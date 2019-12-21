package backup

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/output"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"runtime"
	"sync"
)

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
