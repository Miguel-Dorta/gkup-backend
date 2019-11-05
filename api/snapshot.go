package api

// SnapshotList represents a serializable list of Snapshots objects.
type SnapshotList struct {
	SList []*Snapshots `json:"snapshots"`
}

// Snapshots represents all the snapshots taken under the same name.
type Snapshots struct {
	Name  string  `json:"name"`
	Times []*Time `json:"times"`
}

// Snapshot represents just one snapshot.
type Snapshot struct {
	Name string `json:"name"`
	Time *Time  `json:"time"`
}

// Time represents the time when a snapshot was taken.
type Time struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
