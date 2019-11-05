package check

import (
	"encoding/json"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"github.com/Miguel-Dorta/gkup-backend/pkg/threadSafe"
	"io"
	"time"
)

func statusPrinter(total int, progress *threadSafe.Counter, outWriter io.Writer) {
	seconds := time.NewTicker(time.Second).C
	for range seconds {
		current := progress.Get()

		data, _ := json.Marshal(api.Status{
			Current: current,
			Total:   total,
		})
		_, _ = outWriter.Write(data)

		if current >= total {
			return
		}
	}
}
