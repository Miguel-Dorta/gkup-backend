package check

import (
	"encoding/json"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"io"
	"time"
)

func statusPrinter(outWriter io.Writer) {
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
