package transport

import (
	"RPC_go/discover"
	"RPC_go/transport/connpool"
	"time"
)

type ClientTransportOptions struct {
	Target      string
	ServiceName string
	Network     string
	Pool        connpool.Pool
	Discover    discover.Discover
	Timeout     time.Duration
}

type Option func(o *ClientTransportOptions)

func WithServiceName(name string) Option {
	return func(o *ClientTransportOptions) {
		o.ServiceName = name
	}
}

func WithTarget(target string) Option {
	return func(o *ClientTransportOptions) {
		o.Target = target
	}
}

func WithNetwork(network string) Option {
	return func(o *ClientTransportOptions) {
		o.Network = network
	}
}
