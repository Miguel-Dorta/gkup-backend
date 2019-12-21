package backup

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"sync"
)

type file struct {
	files.File
	RealPath string
}

type safeFileList struct {
	list  []*file
	pos   int
	mutex sync.Mutex
}

func (l *safeFileList) next() *file {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.pos >= len(l.list) {
		return nil
	}
	s := l.list[l.pos]
	l.pos++

	return s
}
