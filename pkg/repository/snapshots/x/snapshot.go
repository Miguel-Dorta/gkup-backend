package x

type Encoder interface {
	Encode(f *File) error
	Close() error
}

type Decoder interface {
	Decode(f *File) error
	More() bool
	Close() error
}

type File struct {
	RelPath string
	Hash []byte
	Size int64
}
