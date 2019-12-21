package threadSafe

import "sync"

type Counter struct {
	c int
	m sync.Mutex
}

func (c *Counter)Add(i int) {
	c.m.Lock()
	c.c += i
	c.m.Unlock()
}

func (c *Counter)Reset() {
	c.m.Lock()
	c.c = 0
	c.m.Unlock()
}

func (c *Counter)Get() int {
	return c.c
}
