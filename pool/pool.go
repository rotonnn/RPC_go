package pool

import (
	"context"
	"errors"
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
		return p.wrapConn(conn)
	}
}

func (p *ConnPool) wrapConn(conn net.Conn) (*Conn, error) {
	resConn := &Conn{
		Conn:          conn,
		dialTimeout:   p.dialTimeout,
		available:     true,
		idleStartTime: time.Now(),
	}
	return resConn, nil
}

func (p *pool) NewConnPool(ctx context.Context, network string, addr string) (*ConnPool, error) {
	return nil, nil
}

func ConnPoolRegister() {

}
