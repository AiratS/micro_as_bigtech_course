package interceptor

import (
	"context"

	rateLimiter "github.com/AiratS/micro_as_bigtech_course/week_8/rate_limiter/internal/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiterInterceptor struct {
	rateLimiter *rateLimiter.TokenBucketLimiter
}

func NewRateLimiterInterceptor(rl *rateLimiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		rateLimiter: rl,
	}
}

func (i *RateLimiterInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if !i.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "to many requests")
	}

	return handler(ctx, req)
}
