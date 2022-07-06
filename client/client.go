package client

import (
	"RPC_go/codec"
	"RPC_go/constant"
	"RPC_go/interceptor"
	"RPC_go/protocol"
	"RPC_go/transport"
	"RPC_go/utils"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
)

type Client interface {
	Invoke(ctx context.Context, req interface{}, resp interface{}, path string, opts ...Option) error
	Call(ctx context.Context, path string, req interface{}, resp interface{}, opts ...Option) error
}

type defaultClient struct {
	opts *Options
}

var DefaultClient = defaultClient{
	opts: &Options{
		protocol: "proto",
	},
}

func (c *defaultClient) Call(ctx context.Context, path string,
	req interface{}, resp interface{}, opts ...Option) error {

	opts = append(opts, WithSerializationType(constant.Proto))
	return c.Invoke(ctx, req, resp, path, opts...)
}

func (c *defaultClient) Invoke(ctx context.Context,
	req interface{}, resp interface{}, path string, opts ...Option) error {
	// Init with optFunc
	for _, opt := range opts {
		opt(c.opts)
	}

	// Set timeout cancel ctx
	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	serviceName, method, err := utils.ParseServicePath(path)
	if err != nil {
		return err
	}
	c.opts.serviceName = serviceName
	c.opts.method = method

	return interceptor.ClientIntercept(ctx, req, resp, c.opts.interceptor, c.doInvoke)
}

func (c *defaultClient) doInvoke(ctx context.Context, req, resp interface{}) error {

	reqBody, err := c.Payload2Byte(ctx, c.opts.protocol, req)
	if err != nil {
		return err
	}

	cliTransport := c.NewClientTransport()
	cliTransportOpt := c.GetClientTransportOptDefault()
	frame, err := cliTransport.Send(ctx, reqBody, cliTransportOpt)
	if err != nil {
		log.Fatalf("[client invoke]Failed to send packet, err=%v", err)
		return err
	}

	resp = c.Byte2Payload(ctx, c.opts.protocol, frame)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultClient) Payload2Byte(ctx context.Context, protocol string, req interface{}) ([]byte, error) {

	serialization := codec.GetSerialization(c.opts.serializationType)
	payload, err := serialization.Serialize(req)
	if err != nil {
		log.Fatalf("[client invoke] Serialize failed, err=%v", err)
		return nil, err
	}

	request := NewRequest(ctx, c, payload)
	reqBuf, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("[client invoke] proto marshal failed, err=%v", err)
		return nil, err
	}

	protoCodec := codec.GetCodec(c.opts.protocol)
	return protoCodec.Encode(reqBuf)
}

func (c *defaultClient) Byte2Payload(ctx context.Context, protoType string, frame []byte) interface{} {
	protoCodec := codec.GetCodec(protoType)
	respBuf, err := protoCodec.Decode(frame)
	if err != nil {
		log.Fatalf("[client invoke] resp buf decode failed, err=%v", err)
		return err
	}

	response := &protocol.Response{}
	err = proto.Unmarshal(respBuf, response)
	if err != nil {
		log.Fatalf("[client invoke] resp buf Unmarshal failed, err=%v", err)
		return err
	}
	if response.ErrCode != 0 {
		return errors.New(fmt.Sprintf("[client invoke] Response not zero, code=%d, err=%s", response.ErrCode, response.ErrTips))
	}

	serialization := codec.GetSerialization(c.opts.serializationType)
	return serialization.Deserialize(response.Payload, response)
}

func (c *defaultClient) NewClientTransport() *transport.ClientTransport {
	return transport.GetClientTransport(c.opts.protocol)
}

func NewRequest(ctx context.Context, c *defaultClient, payload []byte) *protocol.Request {
	path := fmt.Sprintf("/%s/%s", c.opts.serviceName, c.opts.method)
	return &protocol.Request{
		ServicePath: path,
		Payload:     payload,
	}
}

func (c *defaultClient) GetClientTransportOptDefault() []transport.Option {
	return []transport.Option{
		transport.WithTarget(c.opts.target),
		transport.WithServiceName(c.opts.serviceName),
		transport.WithNetwork(c.opts.network),
	}

}
