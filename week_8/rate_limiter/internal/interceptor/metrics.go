package interceptor

import (
	"context"

	"github.com/AiratS/micro_as_bigtech_course/week_8/rate_limiter/internal/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	metric.IncRequestCounter()

	return handler(ctx, req)
}
