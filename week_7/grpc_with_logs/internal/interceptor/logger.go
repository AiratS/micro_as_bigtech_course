package interceptor

import (
	"context"
	"time"

	"github.com/airats/micro_as_bigtech_course/week_7/grpc_with_logs/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}

	logger.Info("request",
		zap.String("method", info.FullMethod),
		zap.Any("req", req),
		zap.Duration("duration", time.Since(now)),
	)

	return res, err
}
