package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// prodLog()
	// sugarLog()
	samplingLog()
}

func prodLog() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Info("Hello world from logs")

	logger.Info("Prod logger",
		zap.String("username", "airat"),
		zap.Int("id", 132),
		zap.String("provider", "test"),
	)
}

func sugarLog() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Sugar().Infow("suggar logger", "name", "oleg", "ID", 45)
}

func samplingLog() {
	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	prodCfg := zap.NewProductionEncoderConfig()
	prodCfg.TimeKey = "timestamp"
	prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	prodCfg.StacktraceKey = "stack"

	jsonEncoder := zapcore.NewJSONEncoder(prodCfg)
	jsonOutCore := zapcore.NewCore(jsonEncoder, stdout, level)

	samplingCore := zapcore.NewSamplerWithOptions(jsonOutCore, time.Second, 3, 0)

	log := zap.New(samplingCore)

	for i := 0; i <= 10; i++ {
		log.Info("hello world")
		log.Warn("warn")
	}
}
