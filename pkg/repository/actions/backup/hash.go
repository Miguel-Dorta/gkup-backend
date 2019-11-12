package backup

import (
	"encoding/hex"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"sync"
)

func hashFilesParallel(list []*copyInfo, hashAlgorithm string, bufferSize, threads int, progress *threadSafe.Counter) error {
	var (
		err error
		safeList = newSafeCopyInfoList(list)
		wg = new(sync.WaitGroup)
	)

	for i:=0; i<threads; i++ {
		wg.Add(1)
		go func() {
			if subErr := fileHasher(safeList, hashAlgorithm, bufferSize, &err, progress); subErr != nil {
				err = subErr
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return err
}

func fileHasher(ciList *safeCopyInfoList, hashAlgorithm string, bufferSize int, upperError *error, progress *threadSafe.Counter) error {
	h := hash.Algorithms[hashAlgorithm]()
	buf := make([]byte, bufferSize)

	for {
		if *upperError != nil {
			return nil
		}

		ci := ciList.next()
		if ci == nil {
			return nil
		}

		hashBytes, err := hash.HashFile(ci.path, h, buf)
		if err != nil {
			return fmt.Errorf("error found hashing file \"%s\": %s", ci.path, err)
		}

		ci.f.Hash = hex.EncodeToString(hashBytes)
		progress.Add(1)
	}
}
