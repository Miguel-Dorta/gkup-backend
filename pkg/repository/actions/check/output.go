package check

import (
	"encoding/json"
	api "github.com/Miguel-Dorta/gkup-backend/api/check"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"io"
	"time"
)

func statusPrinter(total int, progress *threadSafe.Counter, outWriter io.Writer, quit <-chan bool) {
	seconds := time.NewTicker(time.Second).C
	for {
		select {
		case <-seconds:
			printStatus(total, progress.Get(), outWriter)
		case <-quit:
			printStatus(total, progress.Get(), outWriter)
			return
		}
	}
}

func printStatus(total int, progress int, outWriter io.Writer) {
	data, _ := json.Marshal(api.Status{
		ProgressCurrent: progress,
		ProgressTotal:   total,
	})
	_, _ = outWriter.Write(data)
}

