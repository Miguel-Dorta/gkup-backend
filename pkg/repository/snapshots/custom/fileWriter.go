package custom

import (
	"bufio"
	"io"
	"os"
)

type fileWriter struct {
	w *bufio.Writer
	c io.Closer
}

func newFileWriter(f *os.File) *fileWriter {
	return &fileWriter{
		w: bufio.NewWriterSize(f, bufferSize),
		c: f,
	}
}

func (f *fileWriter) Write(p []byte) (int, error) {
	return f.w.Write(p)
}

func (f *fileWriter) Close() error {
	if err := f.w.Flush(); err != nil {
		return err
	}
	return f.c.Close()
}
