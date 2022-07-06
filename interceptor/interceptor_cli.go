package interceptor

import (
	"context"
	"log"
)

type ClientInvoker func(context.Context, interface{}, interface{}) error
type ClientInterceptor func(context.Context, interface{}, interface{}, ClientInvoker) error

func ClientIntercept(ctx context.Context, req, resp interface{}, interceptors []ClientInterceptor, invoker ClientInvoker) error {
	if len(interceptors) <= 0 {
		return invoker(ctx, req, resp)
	}

	for _, interceptor := range interceptors {
		if err := interceptor(ctx, req, resp, invoker); err != nil {
			log.Default()
			return err
		}
	}

	return nil
}
