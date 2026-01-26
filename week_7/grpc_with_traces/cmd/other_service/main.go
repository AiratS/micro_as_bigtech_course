package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/internal/tracing"
	desc "github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/pkg/other_note_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort    = 50062
	serviceName = "other-service"
)

type server struct {
	desc.UnimplementedOtherNoteV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	return &desc.GetResponse{
		Note: &desc.Note{
			Id: req.GetId(),
			Info: &desc.NoteInfo{
				Title:   gofakeit.BeerName(),
				Content: gofakeit.IPv4Address(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	tracing.Init(serviceName)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen port")
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
	))
	reflection.Register(s)

	desc.RegisterOtherNoteV1Server(s, &server{})

	log.Println("other grpc server is running")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to listen port")
	}
}
