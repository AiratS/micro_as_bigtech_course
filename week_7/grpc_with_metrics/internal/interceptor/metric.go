package interceptor

import (
	"context"
	"time"

	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_logs/internal/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	metric.IncRequestCounter()

	timeStart := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(timeStart)

	if err == nil {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.HistogramResponseTimeObserve("success", duration.Seconds())
	} else {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTimeObserve("error", duration.Seconds())
	}

	return res, err
}
