package pool

import (
	"errors"
	"net"
	"sync"
	"time"
)

type Conn struct {
	net.Conn
	dialTimeout   time.Duration
	connPool      ConnPool
	available     bool
	network       string
	addr          string
	idleStartTime time.Time
	mu            *sync.RWMutex
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return
}

func (c *Conn) MarkUnavailable() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.available = false
}

func (c *Conn) Close() error {
	if c.available {
		return errors.New("[conn] conn still available")
	}

	c.mu.RLock()
	if c.available {
		c.mu.Unlock()
		return errors.New("[conn] conn still available")
	}
	c.Conn.Close()
	c.mu.RUnlock()
	return nil
}
