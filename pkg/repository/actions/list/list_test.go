package list_test

import (
	"bytes"
	"encoding/json"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/list"
	"sort"
	"strings"
	"testing"
)

var expectedOutput = api.SnapshotList{
	SList: []*api.Snapshots{
		{
			Name: "",
			Times: []*api.Time{
				{
					Year:   1998,
					Month:  12,
					Day:    28,
					Hour:   1,
					Minute: 23,
					Second: 45,
				},
				{
					Year:   2018,
					Month:  10,
					Day:    8,
					Hour:   16,
					Minute: 43,
					Second: 32,
				},
			},
		},
		{
			Name: "custom_name",
			Times: []*api.Time{
				{
					Year:   1,
					Month:  1,
					Day:    1,
					Hour:   0,
					Minute: 0,
					Second: 0,
				},
				{
					Year:   2099,
					Month:  12,
					Day:    31,
					Hour:   23,
					Minute: 59,
					Second: 59,
				},
			},
		},
		{
			Name: "MyPC",
			Times: []*api.Time{
				{
					Year:   1969,
					Month:  12,
					Day:    31,
					Hour:   23,
					Minute: 59,
					Second: 59,
				},
				{
					Year:   1970,
					Month:  1,
					Day:    1,
					Hour:   0,
					Minute: 0,
					Second: 0,
				},
			},
		},
		{
			Name: "mypc",
			Times: []*api.Time{
				{
					Year:   2038,
					Month:  1,
					Day:    19,
					Hour:   3,
					Minute: 14,
					Second: 7,
				},
				{
					Year:   2038,
					Month:  1,
					Day:    19,
					Hour:   3,
					Minute: 14,
					Second: 8,
				},
			},
		},
	},
}

func TestList(t *testing.T) {
	var outWriter, errWriter = bytes.NewBuffer(make([]byte, 0, 200)), bytes.NewBuffer(nil)
	list.List("testdata", outWriter, errWriter)

	// Check for no errors
	if errs := errWriter.Bytes(); len(errs) != 0 {
		t.Errorf("errors found: %s", string(errs))
	}

	// Unmarshal output
	var actualOutput api.SnapshotList
	if err := json.Unmarshal(outWriter.Bytes(), &actualOutput); err != nil {
		t.Fatalf("cannot unmarsha output \"%v\": %s", outWriter.Bytes(), err)
	}

	// Sort actual output to match expected output
	sort.Slice(actualOutput.SList, func(i, j int) bool {
		iMin := strings.ToLower(actualOutput.SList[i].Name)
		jMin := strings.ToLower(actualOutput.SList[j].Name)

		if iMin != jMin {
			return iMin < jMin
		}
		return actualOutput.SList[i].Name < actualOutput.SList[j].Name
	})
	for i := range actualOutput.SList {
		sort.Slice(actualOutput.SList[i].Times, func(j, k int) bool {
			jTime := actualOutput.SList[i].Times[j]
			kTime := actualOutput.SList[i].Times[k]

			if jTime.Year != kTime.Year {
				return jTime.Year < kTime.Year
			}
			if jTime.Month != kTime.Month {
				return jTime.Month < kTime.Month
			}
			if jTime.Day != kTime.Day {
				return jTime.Day < kTime.Day
			}
			if jTime.Hour != kTime.Hour {
				return jTime.Hour < kTime.Hour
			}
			if jTime.Minute != kTime.Minute {
				return jTime.Minute < kTime.Minute
			}
			return jTime.Second < kTime.Second
		})
	}

	// Compare
	if len(expectedOutput.SList) != len(actualOutput.SList) {
		t.Fatalf("expected size (%d) is not size found (%d) in SList", len(expectedOutput.SList), len(actualOutput.SList))
	}
	for i := range expectedOutput.SList {
		equalsSnap(expectedOutput.SList[i], actualOutput.SList[i], t)
	}
}

func equalsSnap(expectedSnap, actualSnap *api.Snapshots, t *testing.T) {
	if expectedSnap.Name != actualSnap.Name {
		t.Errorf("name expected (%s) is not name found (%s)", expectedSnap.Name, actualSnap.Name)
		return
	}
	if len(expectedSnap.Times) != len(actualSnap.Times) {
		t.Errorf("expected size (%d) is not size found (%d) in SList", len(expectedSnap.Times), len(actualSnap.Times))
		return
	}
	for j := range expectedSnap.Times {
		if !equalsTime(expectedSnap.Times[j], actualSnap.Times[j]) {
			t.Errorf("expected time (%v) don't match found time (%v)", expectedSnap.Times[j], actualSnap.Times[j])
		}
	}
}

func equalsTime(t1, t2 *api.Time) bool {
	return t1.Year == t2.Year && t1.Month == t2.Month && t1.Day == t2.Day && t1.Hour == t2.Hour && t1.Minute == t2.Minute && t1.Second == t2.Second
}
