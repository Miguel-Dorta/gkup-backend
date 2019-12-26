package threadSafe

import (
	"io"
	"sync"
)

type Writer struct {
	w io.Writer
	m sync.Mutex
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.w.Write(p)
}
