package interceptor

import (
	"context"
	"log"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerInterceptor struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreakerInterceptor(cb *gobreaker.CircuitBreaker) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{
		cb: cb,
	}
}

func (c *CircuitBreakerInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	res, err := c.cb.Execute(func() (interface{}, error) {
		r, er := handler(ctx, req)

		log.Println(r, er)

		return r, er
	})

	if err != nil {
		if err == gobreaker.ErrOpenState {
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}

		return nil, err
	}

	return res, err
}
