package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_logs/internal/interceptor"
	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_logs/internal/metric"
	desc "github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_logs/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brianvoe/gofakeit"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const grpcPort = 50061

type server struct {
	desc.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
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
	ctx := context.Background()
	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("Failed to init metric: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen port: %v", err)
	}

	ser := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.MetricsInterceptor,
	))
	reflection.Register(ser)
	desc.RegisterNoteV1Server(ser, &server{})

	go func() {
		err := runPrometheus()
		if err != nil {
			log.Fatalf("Failed to run prometheus: %v", err)
		}
	}()

	log.Println("running grpc server")
	if err = ser.Serve(lis); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}

	log.Println("Prometheus is running")

	if err := prometheusServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
