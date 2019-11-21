package x

type Writer interface {
	Write(f *File) error
	Close() error
}

type Reader interface {
	ReadNext() (*File, error)
	More() bool
	Close() error
}

type File struct {
	RelPath string
	Hash []byte
	Size int64
}
