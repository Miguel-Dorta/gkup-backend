package backup

import (
	"encoding/json"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"io"
	"sync"
	"time"
)

func startStatus(progress *threadSafe.Counter, total int, outWriter io.Writer) func() {
	quitChan := make(chan bool, 10) // Large buffer so it can be called multiple times
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		statusPrinter(progress, total, outWriter, quitChan)
		wg.Done()
	}()

	return func() {
		quitChan <- true
		wg.Wait()
	}
}

func statusPrinter(current *threadSafe.Counter, total int, outWriter io.Writer, quit <-chan bool) {
	seconds := time.NewTicker(time.Second).C
	for {
		select {
		case <-quit:
			printStatus(current.Get(), total, outWriter)
			return
		case <-seconds:
			progress := current.Get()
			printStatus(progress, total, outWriter)

			if progress >= total {
				return
			}
		}
	}
}

func printStatus(current int, total int, outWriter io.Writer) {
	data, _ := json.Marshal(api.Status{
		Current: current,
		Total:   total,
	})
	_, _ = outWriter.Write(data)
}
