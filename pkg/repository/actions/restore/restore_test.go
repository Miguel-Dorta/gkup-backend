package restore_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/api"
	"github.com/Miguel-Dorta/gkup-backend/pkg/hash"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/actions/restore"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

type expectedFile struct {
	name  string
	hash string
	isDir bool
}

type expectedDir struct {
	path string
	list []expectedFile
}

const (
	repoPath   = "testdata"
	bufferSize = 128 * 1024
)

var testPath = filepath.Join(os.TempDir(), "gkup_restore_test")

var expectedList = []expectedDir{
	{
		path: ".",
		list: []expectedFile{
			{name: "dir0", isDir: true},
			{name: "dir5", isDir: true},
			{name: "file16", isDir: false, hash: "87ecb1828e77509486215cf1d9cb4662ba5dc6e323ca6bef00eb071a41ffc953"},
			{name: "file17", isDir: false, hash: "9ee2e1004b50d1351fb9a06b3a7b529372442cecbc984540704934c05431edad"},
		},
	},
	{
		path: "dir0",
		list: []expectedFile{
			{name: "dir1", isDir: true},
			{name: "dir2", isDir: true},
			{name: "file08", isDir: false, hash: "a78fca73d3cf76c693af93c01d4ca3a2b810fd05e20aafa045a245d610069433"},
			{name: "file09", isDir: false, hash: "c947655ebcbe162c46b9b31040fd66dd75cc98bbb8a8cb4f429b57e79abba319"},
		},
	},
	{
		path: filepath.Join("dir0", "dir1"),
		list: []expectedFile{
			{name: "file00", isDir: false, hash: "022fe12f0d72275b5400121bcc82792452db161e875b371fdc269a8b73b137e0"},
			{name: "file01", isDir: false, hash: "0ed36861e667616827c4468f8933e67d231ba374ccaa201c62879867d71671f2"},
		},
	},
	{
		path: filepath.Join("dir0", "dir2"),
		list: []expectedFile{
			{name: "dir3", isDir: true},
			{name: "file06", isDir: false, hash: "89b00c75f4a3941a1c62b50e67bc037913637ebdddbe86edbcb03feee212d9ab"},
			{name: "file07", isDir: false, hash: "96201c014bf3a116980528b4a6c08804b698431f338ff8af2d54fbc2021cf56a"},
		},
	},
	{
		path: filepath.Join("dir0", "dir2", "dir3"),
		list: []expectedFile{
			{name: "dir4", isDir: true},
			{name: "file04", isDir: false, hash: "4a6f44fbe6d97db2ec6e5562a34c0ba8449ab21f0164e840979b34fc15b661fc"},
			{name: "file05", isDir: false, hash: "80a8af36df93c96314b8f6c8d3d6a33bad33b50a0395418b2b38517b012370b7"},
		},
	},
	{
		path: filepath.Join("dir0", "dir2", "dir3", "dir4"),
		list: []expectedFile{
			{name: "file02", isDir: false, hash: "178ee7e5a462052981e72e450e0f5bc43f61ed183f331317bb25410a2d249381"},
			{name: "file03", isDir: false, hash: "35c5480f9d81aa7281e559b7134bde9cf69573e5a967da20e12b34b54ea63f8b"},
		},
	},
	{
		path: "dir5",
		list: []expectedFile{
			{name: "dir6", isDir: true},
			{name: "dir7", isDir: true},
			{name: "file14", isDir: false, hash: "ea90ea27825283a2b3c2b6580e2390f0dc5d3a9af9ccce7637ccf59ccb1c5b13"},
			{name: "file15", isDir: false, hash: "ff92b70a66ee368fbfc39b0d80245bb29d9d4e0da4abc8cd7244e26a4bb10842"},
		},
	},
	{
		path: filepath.Join("dir5", "dir6"),
		list: []expectedFile{
			{name: "file10", isDir: false, hash: "ccbab00ec41e1a0a57571bc7689e94b35dda8f056b4d24c9d8bcd4ba68b7d697"},
			{name: "file11", isDir: false, hash: "ccfd79d83959bb758db3a771100ba9e15e249b2e7348f3ae574fc484bc342241"},
		},
	},
	{
		path: filepath.Join("dir5", "dir7"),
		list: []expectedFile{
			{name: "file12", isDir: false, hash: "d828e2d9650e9999b8e3ff0bbccb8f59f7fbb28fbe7afdf18270bc175aecbd8a"},
			{name: "file13", isDir: false, hash: "e062a68f09c9103affd279f9c181edcb81db476dbda64d49ef8d38ce9f5a6cb0"},
		},
	},
}

func TestRestore(t *testing.T) {
	if err := os.Mkdir(testPath, 0777); err != nil {
		t.Fatalf("cannot create tmp dir (%s): %s", testPath, err)
	}
	defer os.RemoveAll(testPath)

	var (
		out []byte
		errs []string
		err error
		list []os.FileInfo
	)

	// Non-existent destiny test
	out, errs, err, list = nil, nil, nil, nil
	out, errs = execTest(filepath.Join(testPath, "nonexistent"), &api.Snapshot{
		Name: "",
		Time: &api.Time{
			Year:   2019,
			Month:  11,
			Day:    5,
			Hour:   13,
			Minute: 28,
			Second: 0,
		},
	})
	checkRestore(filepath.Join(testPath, "nonexistent"), t)
	checkOut(out, t)
	printErrs("non-existent destiny", errs, t)

	// Empty destiny test
	out, errs, err, list = nil, nil, nil, nil
	if err := os.Mkdir(filepath.Join(testPath, "empty"), 0777); err != nil {
		t.Errorf("cannot create dir")
	}
	out, errs = execTest(filepath.Join(testPath, "empty"), &api.Snapshot{
		Name: "snapWithName",
		Time: &api.Time{
			Year:   2010,
			Month:  12,
			Day:    31,
			Hour:   23,
			Minute: 59,
			Second: 59,
		},
	})
	checkRestore(filepath.Join(testPath, "empty"), t)
	checkOut(out, t)
	printErrs("empty destiny", errs, t)

	// Non-empty destiny test
	out, errs, err, list = nil, nil, nil, nil
	if err := os.Mkdir(filepath.Join(testPath, "nonempty"), 0777); err != nil {
		t.Errorf("cannot create dir")
	}
	if err := createEmptyFile(filepath.Join(testPath, "nonempty", "file.txt")); err != nil {
		t.Error("cannot create empty file")
	}
	out, errs = execTest(filepath.Join(testPath, "nonempty"), &api.Snapshot{
		Name: "snapWithName",
		Time: &api.Time{
			Year:   2010,
			Month:  12,
			Day:    31,
			Hour:   23,
			Minute: 59,
			Second: 59,
		},
	})
	if len(errs) == 0 {
		t.Error("no error was printed in nonempty destiny")
	}
	if len(out) != 0 {
		t.Errorf("status was printed in nonempty destiny: %s", string(out))
	}
	list, err = utils.ListDir(filepath.Join(testPath, "nonempty"))
	if err != nil {
		t.Error("cannot list nonempty dir")
	}
	if len(list) != 1 {
		t.Errorf("nonempty dir was modified (current number of files: %d)", len(list))
	}

	// Non-existing snapshot test
	out, errs, err, list = nil, nil, nil, nil
	out, errs = execTest(filepath.Join(testPath, "nonexistingSnap"), &api.Snapshot{
		Name: "",
		Time: &api.Time{
			Year:   2000,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 0,
			Second: 0,
		},
	})
	if len(errs) == 0 {
		t.Error("no error was printed in non-existing snapshot")
	}
	if len(out) != 0 {
		t.Errorf("status was printed in non-existing snapshot: %s", string(out))
	}

	// Non-existing file to copy test
	out, errs, err, list = nil, nil, nil, nil
	out, errs = execTest(filepath.Join(testPath, "nonexistingFile"), &api.Snapshot{
		Name: "",
		Time: &api.Time{
			Year:   2019,
			Month:  11,
			Day:    5,
			Hour:   13,
			Minute: 28,
			Second: 99,
		},
	})
	if len(errs) == 0 {
		t.Error("no error was printed in non-existing snapshot")
	}
	checkOut(out, t)

	// TODO symlink
}

func execTest(restorePath string, restoreSnap *api.Snapshot) ([]byte, []string) {
	outWriter, errWriter := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	restore.Restore(repoPath, restorePath, bufferSize, restoreSnap, outWriter, errWriter)

	parts := bytes.Split(errWriter.Bytes(), []byte{'\n'})
	errs := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		errs = append(errs, string(part))
	}

	return outWriter.Bytes(), errs
}

func printErrs(testname string, errs []string, t *testing.T) {
	if len(errs) == 0 {
		return
	}

	t.Errorf("errors found in %s test", testname)
	for _, err := range errs {
		t.Error(err)
	}
}

func checkOut(data []byte, t *testing.T) {
	if len(data) == 0 {
		t.Error("no output was found")
		return
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	for dec.More() {
		var status api.Status
		if err := dec.Decode(&status); err != nil {
			t.Errorf("error decoding status (%s): %s", string(data), err)
			continue
		}

		if status.Current < 0 || status.Current > status.Total || status.Total != 18 {
			t.Errorf("incorrect status: current (%d) - total (%d)", status.Current, status.Total)
		}
	}
}

func checkRestore(restorePath string, t *testing.T) {
	for _, dir := range expectedList {
		if err := checkPath(filepath.Join(restorePath, dir.path), dir.list); err != nil {
			t.Errorf("error checking restoration in path \"%s\": %s", filepath.Join(restorePath, dir.path), err)
		}
	}
}

func checkPath(path string, l []expectedFile) error {
	list, err := utils.ListDir(path)
	if err != nil {
		return fmt.Errorf("cannot list path (%s): %s", path, err)
	}

	actualFiles := make([]expectedFile, len(list))
	for i, fi := range list {
		f := expectedFile{
			name:  fi.Name(),
			isDir: fi.IsDir(),
		}
		if !fi.IsDir() {
			h, err := hash.HashFile(filepath.Join(path, fi.Name()), hash.Algorithms["sha256"](), make([]byte, bufferSize))
			if err != nil {
				return fmt.Errorf("cannot hash file (%s): %s", filepath.Join(path, fi.Name()), err)
			}
			f.hash = hex.EncodeToString(h)
		}
		actualFiles[i] = f
	}

	sort.Slice(actualFiles, func(i, j int) bool {
		return actualFiles[i].name < actualFiles[j].name
	})

	return checkListFileInfo(actualFiles, l)
}

func checkListFileInfo(actual []expectedFile, expected []expectedFile) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("number of expected files (%d) is not equal to the number of files found (%d)", len(expected), len(actual))
	}

	for i := range actual {
		expectedF := expected[i]
		actualF := actual[i]

		if actualF.name != expectedF.name || actualF.isDir != expectedF.isDir {
			return fmt.Errorf("expected file (name: %s, isDir: %t) is not equal to file found (name: %s, isDir: %t)", expectedF.name, expectedF.isDir, actualF.name, actualF.isDir)
		}
		if !expectedF.isDir && expectedF.hash != actualF.hash {
			return fmt.Errorf("hashes don't match: expected (%s) - found (%s)", expectedF.hash, actualF.hash)
		}
	}
	return nil
}

func createEmptyFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}
