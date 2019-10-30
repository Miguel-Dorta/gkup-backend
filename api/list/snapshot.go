package list

type SnapshotList struct {
	SList []*Snapshot `json:"snapshots"`
}

type Snapshot struct {
	Name  string  `json:"name"`
	Times []*Time `json:"times"`
}

type Time struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
