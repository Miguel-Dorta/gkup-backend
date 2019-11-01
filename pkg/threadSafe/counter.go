package threadSafe

import "sync"

type Counter struct {
	c int
	m sync.Mutex
}

func (c *Counter)Add(i int) {
	c.m.Lock()
	defer c.m.Unlock()
	c.c += i
}

func (c *Counter)Get() int {
	return c.c
}
