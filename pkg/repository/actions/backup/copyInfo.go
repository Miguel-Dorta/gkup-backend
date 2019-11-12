package backup

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/snapshots"
	"sync"
)

type copyInfo struct {
	f    *snapshots.File
	path string
}

type safeCopyInfoList struct {
	list []*copyInfo
	index int
	m sync.Mutex
}

func newSafeCopyInfoList(l []*copyInfo) *safeCopyInfoList {
	return &safeCopyInfoList{list: l}
}

func (l *safeCopyInfoList) next() *copyInfo {
	l.m.Lock()
	defer l.m.Unlock()

	if l.index >= len(l.list) {
		return nil
	}
	l.index++
	return l.list[l.index-1]
}
