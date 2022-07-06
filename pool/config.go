package pool

import "time"

type Options struct {
	maxCap      int
	initialCap  int
	maxIdle     int
	idleTimeout time.Duration
	dialTimeout time.Duration
}

type Option func(options *Options)

func WithMaxCap(cap int) Option {
	return func(o *Options) {
		o.maxCap = cap
	}
}
func WithIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = t
	}
}
func WithDialTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.dialTimeout = t
	}
}
