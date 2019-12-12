package custom

import (
	"bufio"
	"io"
	"os"
)

type fileReader struct {
	r *bufio.Reader
	c io.Closer
}

func newFileReader(f *os.File) *fileReader {
	return &fileReader{
		r: bufio.NewReaderSize(f, bufferSize),
		c: f,
	}
}

func (f *fileReader) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f *fileReader) Close() error {
	return f.c.Close()
}
