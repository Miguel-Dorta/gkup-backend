package create

import (
	"fmt"
	"io"
)

type logger struct {
	io.Writer
}

func (l *logger) errorf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(l.Writer, format + "\n", a)
}
