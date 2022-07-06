package transport

import (
	"context"
)

var clientTransportMap = map[string]*ClientTransport{
	"proto": &ClientTransport{},
}

type ClientTransport interface {
	Send(ctx context.Context, data []byte, option []Option) ([]byte, error)
}

func GetClientTransport(protoType string) *ClientTransport {
	if t, ok := clientTransportMap[protoType]; ok && t != nil {
		return t
	}
	return &ClientTransport{}
}

type ClientTransport struct {
	opts *ClientTransportOptions
}

func (t *ClientTransport) Send(ctx context.Context, data []byte, option []Option) ([]byte, error) {
	return nil, nil
}
