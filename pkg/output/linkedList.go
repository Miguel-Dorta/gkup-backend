package output

type linkedList struct {
	first, last *node
}

type node struct {
	value *statusInfo
	next  *node
}

func (l *linkedList) Push(value statusInfo) {
	newNode := &node{value: &value}

	if l.first == nil {
		l.first = newNode
		l.last = newNode
		return
	}

	l.last.next = newNode
	l.last = newNode
}

func (l *linkedList) PopAndReset() (n *node) {
	n = l.first
	l.first = nil
	return
}
