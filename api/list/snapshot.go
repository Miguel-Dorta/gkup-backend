package list

// SnapshotList represents a serializable list of Snapshot objects.
type SnapshotList struct {
	SList []*Snapshot `json:"snapshots"`
}

// Snapshot represents all the snapshots taken under the same name.
type Snapshot struct {
	Name  string  `json:"name"`
	Times []*Time `json:"times"`
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
