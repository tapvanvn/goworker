package goworker

import "sync"

type counter struct {
	mu    sync.Mutex
	Value int
}

func (c *counter) Inc() {
	c.mu.Lock()
	c.Value++
	c.mu.Unlock()
}

func (c *counter) Desc() {
	c.mu.Lock()
	c.Value--
	c.mu.Unlock()
}

func (c *counter) UnsafeLock() {
	c.mu.Lock()
}

func (c *counter) UnsafeUnlock() {
	c.mu.Unlock()
}
