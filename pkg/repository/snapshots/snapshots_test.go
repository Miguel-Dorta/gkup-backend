package snapshots_test

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"os"
	"path/filepath"
	"testing"
)

var valid = &snapshots.Snapshot{
	Version: "v1.0.0",
	Files: snapshots.Directory{
		Name:  ".",
		Dirs:  []*snapshots.Directory{
			{
				Name: "dir0",
				Dirs:  []*snapshots.Directory{
					{
						Name: "dir1",
						Dirs:  nil,
						Files: []*snapshots.File{
							{
								Name: "file00",
								Hash: "022fe12f0d72275b5400121bcc82792452db161e875b371fdc269a8b73b137e0",
								Size: 5120,
							},
							{
								Name: "file01",
								Hash: "0ed36861e667616827c4468f8933e67d231ba374ccaa201c62879867d71671f2",
								Size: 5120,
							},
						},
					},
					{
						Name: "dir2",
						Dirs:  []*snapshots.Directory{
							{
								Name: "dir3",
								Dirs:  []*snapshots.Directory{
									{
										Name: "dir4",
										Dirs:  nil,
										Files: []*snapshots.File{
											{
												Name: "file02",
												Hash: "178ee7e5a462052981e72e450e0f5bc43f61ed183f331317bb25410a2d249381",
												Size: 5120,
											},
											{
												Name: "file03",
												Hash: "35c5480f9d81aa7281e559b7134bde9cf69573e5a967da20e12b34b54ea63f8b",
												Size: 5120,
											},
										},
									},
								},
								Files: []*snapshots.File{
									{
										Name: "file04",
										Hash: "4a6f44fbe6d97db2ec6e5562a34c0ba8449ab21f0164e840979b34fc15b661fc",
										Size: 5120,
									},
									{
										Name: "file05",
										Hash: "80a8af36df93c96314b8f6c8d3d6a33bad33b50a0395418b2b38517b012370b7",
										Size: 5120,
									},
								},
							},
						},
						Files: []*snapshots.File{
							{
								Name: "file06",
								Hash: "89b00c75f4a3941a1c62b50e67bc037913637ebdddbe86edbcb03feee212d9ab",
								Size: 5120,
							},
							{
								Name: "file07",
								Hash: "96201c014bf3a116980528b4a6c08804b698431f338ff8af2d54fbc2021cf56a",
								Size: 5120,
							},
						},
					},
				},
				Files: []*snapshots.File{
					{
						Name: "file08",
						Hash: "a78fca73d3cf76c693af93c01d4ca3a2b810fd05e20aafa045a245d610069433",
						Size: 5120,
					},
					{
						Name: "file09",
						Hash: "c947655ebcbe162c46b9b31040fd66dd75cc98bbb8a8cb4f429b57e79abba319",
						Size: 5120,
					},
				},
			},
			{
				Name: "dir5",
				Dirs:  []*snapshots.Directory{
					{
						Name: "dir6",
						Dirs:  nil,
						Files: []*snapshots.File{
							{
								Name: "file10",
								Hash: "ccbab00ec41e1a0a57571bc7689e94b35dda8f056b4d24c9d8bcd4ba68b7d697",
								Size: 5120,
							},
							{
								Name: "file11",
								Hash: "ccfd79d83959bb758db3a771100ba9e15e249b2e7348f3ae574fc484bc342241",
								Size: 5120,
							},
						},
					},
					{
						Name: "dir7",
						Dirs:  nil,
						Files: []*snapshots.File{
							{
								Name: "file12",
								Hash: "d828e2d9650e9999b8e3ff0bbccb8f59f7fbb28fbe7afdf18270bc175aecbd8a",
								Size: 5120,
							},
							{
								Name: "file13",
								Hash: "e062a68f09c9103affd279f9c181edcb81db476dbda64d49ef8d38ce9f5a6cb0",
								Size: 5120,
							},
						},
					},
				},
				Files: []*snapshots.File{
					{
						Name: "file14",
						Hash: "ea90ea27825283a2b3c2b6580e2390f0dc5d3a9af9ccce7637ccf59ccb1c5b13",
						Size: 5120,
					},
					{
						Name: "file15",
						Hash: "ff92b70a66ee368fbfc39b0d80245bb29d9d4e0da4abc8cd7244e26a4bb10842",
						Size: 5120,
					},
				},
			},
		},
		Files: []*snapshots.File{
			{
				Name: "file16",
				Hash: "87ecb1828e77509486215cf1d9cb4662ba5dc6e323ca6bef00eb071a41ffc953",
				Size: 5120,
			},
			{
				Name: "file17",
				Hash: "9ee2e1004b50d1351fb9a06b3a7b529372442cecbc984540704934c05431edad",
				Size: 5120,
			},
		},
	},
}

func TestRead(t *testing.T) {
	// check valid
	snap, err := snapshots.Read(filepath.Join("testdata", "valid.json"))
	if err != nil {
		t.Errorf("error found in valid test: %s", err)
	}
	checkSnap(snap, valid, t)

	// check invalid
	if _, err := snapshots.Read(filepath.Join("testdata", "invalid.json")); err == nil {
		t.Error("no error found in invalid test")
	}

	// check empty
	if _, err := snapshots.Read(filepath.Join("testdata", "empty")); err == nil {
		t.Error("no error found in empty test")
	}

	// check nonexistent
	if _, err := snapshots.Read(filepath.Join("testdata", "nonexistent")); err == nil {
		t.Error("no error found in nonexistent test")
	}
}

func TestWrite(t *testing.T) {
	tmpPath := filepath.Join(os.TempDir(), "gkup_snapshots_write_test.json")
	defer os.Remove(tmpPath)

	if err := snapshots.Write(tmpPath, valid); err != nil {
		t.Fatalf("error writing file (%s): %s", tmpPath, err)
	}

	snap, err := snapshots.Read(tmpPath)
	if err != nil {
		t.Fatalf("error reading file (%s): %s", tmpPath, err)
	}

	checkSnap(snap, valid, t)
}

func checkSnap(s, expected *snapshots.Snapshot, t *testing.T) {
	if s.Version != expected.Version {
		t.Errorf("version found (%s) is not what expected (%s)", s.Version, expected.Version)
	}
	checkDir(&s.Files, &expected.Files, t)
}

func checkDir(d, expected *snapshots.Directory, t *testing.T) {
	if d.Name != expected.Name {
		t.Errorf("name found (%s) was not expected (%s)", d.Name, expected.Name)
		return
	}

	// Check files
	if len(d.Files) != len(expected.Files) {
		t.Errorf("number of files in dir %s (%d) is not expected (%d)", d.Name, len(d.Files), len(expected.Files))
		return
	}
	for i := range d.Files {
		checkFile(d.Files[i], expected.Files[i], t)
	}

	// Check dirs
	if len(d.Dirs) != len(expected.Dirs) {
		t.Errorf("number of dirs in dir %s (%d) is not expected (%d)", d.Name, len(d.Dirs), len(expected.Dirs))
		return
	}
	for i := range d.Dirs {
		checkDir(d.Dirs[i], expected.Dirs[i], t)
	}
}

func checkFile(f, expected *snapshots.File, t *testing.T) {
	if f.Name != expected.Name {
		t.Errorf("name found (%s) was not expected (%s)", f.Name, expected.Name)
		return
	}
	if f.Hash != expected.Hash {
		t.Errorf("hash found (%s) was not expected (%s) in file %s", f.Hash, expected.Hash, f.Name)
		return
	}
	if f.Size != expected.Size {
		t.Errorf("size found (%d) was not expected (%d) in file %s", f.Size, expected.Size, f.Name)
		return
	}
}
