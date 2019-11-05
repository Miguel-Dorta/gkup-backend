package restore

import (
	"encoding/json"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"io"
	"time"
)

func statusPrinter(total int, progress *int, outWriter io.Writer) {
	seconds := time.NewTicker(time.Second).C
	for range seconds {
		current := *progress

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
