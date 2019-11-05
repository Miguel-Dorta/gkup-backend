package check_test

import (
	"bytes"
	"encoding/json"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/check"
	"io"
	"runtime"
	"strings"
	"testing"
)

var expectedErrs = map[string]bool{
	// Invalid name
	"error checking file \"testdata/files/96/random-file\": error parsing filename (testdata/files/96/random-file): error decoding hash from name: encoding/hex: invalid byte: U+0072 'r'": false,
	// Modified content
	"error checking file \"testdata/files/64/64010f4b56691585c72c27067a86cd4e447c9c73fe8218e75d92d503ea05d6ad-5120\": hash found (a1f6080b35cb838380f9878a0e0d6f88c96699ba117760f6490f6931b058486f) doesn't match expected hash (64010f4b56691585c72c27067a86cd4e447c9c73fe8218e75d92d503ea05d6ad)": false,
	// Modified filename (hash)
	"error checking file \"testdata/files/a7/a78fca73d3cf76c693af93c01d4ca3a2b810fd05e20aafa045a235d610069433-5120\": hash found (a78fca73d3cf76c693af93c01d4ca3a2b810fd05e20aafa045a245d610069433) doesn't match expected hash (a78fca73d3cf76c693af93c01d4ca3a2b810fd05e20aafa045a235d610069433)": false,
	// Modified content
	"error checking file \"testdata/files/98/98d9c7552d34d9a179b919836555a408db32451fc9d50eaa86a38428dfd7bd9c-5120\": size found (5122) doesn't match expected size (5120)": false,
	// Modified filename (size)
	"error checking file \"testdata/files/cc/ccfd79d83959bb758db3a771100ba9e15e249b2e7348f3ae574fc484bc342241-5121\": size found (5120) doesn't match expected size (5121)": false,
}

func TestCheck(t *testing.T) {
	var outWriter, errWriter = bytes.NewBuffer(make([]byte, 0, 1024)), bytes.NewBuffer(make([]byte, 0, 1024))
	check.Check("testdata", runtime.NumCPU(), 128*1024, outWriter, errWriter)

	checkErrors(errWriter.String(), t)
	checkStatus(outWriter, t)
}

func checkErrors(errsStr string, t *testing.T) {
	errs := strings.Split(errsStr, "\n")
	for _, err := range errs {
		if len(err) == 0 {
			continue
		}

		if _, exists := expectedErrs[err]; !exists {
			t.Errorf("unexpected errorTXT: %s", err)
			continue
		}
		expectedErrs[err] = true
	}

	for k, v := range expectedErrs {
		if !v {
			t.Errorf("not catched error: %s", k)
		}
	}
}

func checkStatus(statusReader io.Reader, t *testing.T) {
	dec := json.NewDecoder(statusReader)
	counter := 0

	for dec.More() {
		counter++
		var statusJSON api.Status

		if err := dec.Decode(&statusJSON); err != nil {
			t.Errorf("error decoding status: %s", err)
			continue
		}

		if statusJSON.Total < statusJSON.Current {
			t.Errorf("error in status: progress (%d) is greater than total (%d)", statusJSON.Current, statusJSON.Total)
			continue
		}
	}

	if counter == 0 {
		t.Error("not JSON status found!")
	}
}
