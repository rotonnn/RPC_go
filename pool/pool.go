package pool

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

type Pool interface {
	Get()
}

type pool struct {
	options     *Options
	addrPoolMap *sync.Map
}

func (p *pool) NewConnPool(ctx context.Context, network string, addr string) (*ConnPool, error) {
	connPool := &ConnPool{
		initialCap:  p.options.initialCap,
		maxCap:      p.options.maxCap,
		idleTimeout: p.options.idleTimeout,
		dialTimeout: p.options.dialTimeout,
		mu:          &sync.RWMutex{},
		connChan:    make(chan *Conn, p.options.maxCap),
		Dial: func(ctx context.Context) (net.Conn, error) {
			var conn net.Conn
			select {
			case conn = <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			timeout := p.options.dialTimeout
			if t, ok := ctx.Deadline(); ok {
				timeout = t.Sub(time.Now())
			}

			return net.DialTimeout(network, addr, timeout)
		},
	}

	for i := 0; i < p.options.initialCap; i += 1 {
		rawConn, err := connPool.Dial(ctx)
		if err != nil {
			continue
		}
		connPool.Put(connPool.wrapConn(rawConn))
	}

	return connPool, nil
}

func (p *pool) ConnPoolRegister(addr string, cp *ConnPool) {
	p.addrPoolMap.Store(addr, cp)
}

func (p *pool) Get(ctx context.Context, network string, addr string) (net.Conn, error) {
	var (
		connPool *ConnPool
		err      error
	)

	if cp, ok := p.addrPoolMap.Load(addr); ok {
		if connPool, ok := cp.(ConnPool); ok {
			return connPool.Get(ctx)
		}
	}

	connPool, err = p.NewConnPool(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	p.addrPoolMap.Store(addr, connPool)
	return connPool.Get(ctx)
}

type ConnPool struct {
	net.Conn
	initialCap  int
	maxCap      int
	idleTimeout time.Duration
	dialTimeout time.Duration
	mu          *sync.RWMutex
	connChan    chan *Conn

	Dial func(ctx context.Context) (net.Conn, error)
}

func (p *ConnPool) Get(ctx context.Context) (net.Conn, error) {
	if p.connChan == nil {
		return nil, errors.New("connChan is nil")
	}

	select {
	case conn := <-p.connChan:
		if conn == nil || !conn.available {
			return nil, errors.New("[conn] conn closed")
		}
		return conn, nil
	default:
		conn, err := p.Dial(ctx)
		if err != nil {
			return nil, err
		}
		return p.wrapConn(conn), nil
	}
}

func (p *ConnPool) Put(c *Conn) error {
	if c == nil {
		return errors.New("connection closed")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if p.connChan == nil {
		c.MarkUnavailable()
		c.Close()
		return errors.New("connPool chan closed")
	}

	select {
	case p.connChan <- c:
		return nil
	default:
		return c.Close()
	}
}

func (p *ConnPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	conns := p.connChan
	p.connChan = nil

	if conns == nil {
		return nil
	}
	close(conns)
	for conn := range conns {
		conn.MarkUnavailable()
		conn.Close()
	}

	return nil
}

func (p *ConnPool) wrapConn(conn net.Conn) *Conn {
	resConn := &Conn{
		Conn:          conn,
		dialTimeout:   p.dialTimeout,
		available:     true,
		idleStartTime: time.Now(),
	}
	return resConn
}

func (p *ConnPool) registerChecker(internal time.Duration, checker func(conn *Conn) bool) {
	if internal <= 0 || checker == nil {
		return
	}

	go func() {
		for {
			time.Sleep(internal)

			for i := 0; i < len(p.connChan); i += 1 {
				select {
				case conn := <-p.connChan:
					if !checker(conn) {
						conn.MarkUnavailable()
						conn.Close()
						break
					} else {
						p.Put(conn)
					}
				default:
					break
				}
			}
		}
	}()
}

func (p *ConnPool) Checker(conn *Conn) bool {
	if conn.idleStartTime.Add(p.idleTimeout).Before(time.Now()) || !isConnAlive() {
		return false
	}
	return true
}

func isConnAlive(conn *Conn) bool {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	oneByte := make([]byte, 1)
	if n, err := conn.Read(oneByte); n > 0 || err == io.EOF {
		return false
	}

	conn.SetReadDeadline(time.Time{})
	return true
}
