package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AiratS/micro_as_bigtech_course/week_8/rate_limiter/internal/interceptor"
	"github.com/AiratS/micro_as_bigtech_course/week_8/rate_limiter/internal/rate_limiter"
	desc "github.com/AiratS/micro_as_bigtech_course/week_8/rate_limiter/pkg/note_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

const grpcPort = 50061

type server struct {
	desc.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Note id: %d", req.GetId())

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}

	rateLimiter := rate_limiter.NewTokenBucketLimiter(ctx, 10, time.Second)

	s := grpc.NewServer(grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			interceptor.NewRateLimiterInterceptor(rateLimiter).Unary,
		),
	))
	reflection.Register(s)

	desc.RegisterNoteV1Server(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
