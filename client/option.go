package client

import (
	"RPC_go/interceptor"
	"RPC_go/transport"
	"time"
)

type Options struct {
	serviceName string        `json:"service_name"`
	method      string        `json:"method"`
	target      string        `json:"target"`
	network     string        `json:"network"`
	timeout     time.Duration `json:"timeout"`

	protocol          string //todo change to enum in future
	serializationType string //todo change to enum in future

	transportOpts transport.ClientTransportOptions
	interceptor   []interceptor.ClientInterceptor
	discoverName  string //todo change to enum in future
}

type Option func(*Options)

func WithServiceName(name string) Option {
	return func(o *Options) {
		o.serviceName = name
	}
}
func WithMethod(method string) Option {
	return func(o *Options) {
		o.method = method
	}
}

func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

func WithNetworkl(network string) Option {
	return func(o *Options) {
		o.network = network
	}
}

func WithSerializationType(serializationType string) Option {
	return func(o *Options) {
		o.serializationType = serializationType
	}
}
