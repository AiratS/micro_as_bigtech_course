package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/airats/micro_as_bigtech_course/week_7/grpc_with_logs/internal/interceptor"
	"github.com/airats/micro_as_bigtech_course/week_7/grpc_with_logs/internal/logger"
	desc "github.com/airats/micro_as_bigtech_course/week_7/grpc_with_logs/pkg/note_v1"
	"github.com/brianvoe/gofakeit"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var logLevel = flag.String("l", "info", "Log level")

const grpcPort = 50061

type server struct {
	desc.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("get", zap.Int64("ID", req.GetId()))

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: req.GetId(),
			Info: &desc.NoteInfo{
				Title:    gofakeit.BeerName(),
				Content:  gofakeit.IPv4Address(),
				Author:   gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen port: %d", grpcPort)
	}

	logger.Init(getZapCore(getAtomicLevel()))

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			interceptor.LoggerInterceptor,
		),
	))

	reflection.Register(grpcServer)

	desc.RegisterNoteV1Server(grpcServer, &server{})

	log.Println("Running server")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func getZapCore(level zap.AtomicLevel) zapcore.Core {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 4,
		MaxAge:     7,
	})
	prodCfg := zap.NewProductionEncoderConfig()
	prodCfg.TimeKey = "timestamp"
	prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(prodCfg)

	stdout := zapcore.AddSync(os.Stdout)
	devCfg := zap.NewDevelopmentEncoderConfig()
	devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(devCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
