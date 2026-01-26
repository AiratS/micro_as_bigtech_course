package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/internal/client/rpc"
	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/internal/client/rpc/other_service"
	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/internal/interceptor"
	"github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/internal/tracing"
	desc "github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/pkg/note_v1"
	descOther "github.com/AiratS/micro_as_bigtech_course/week_7/grpc_with_traces/pkg/other_note_v1"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort      = 50061
	grpcOtherPort = 50062
	serviceName   = "test-service"
)

type server struct {
	desc.UnimplementedNoteV1Server

	otherServiceClient rpc.OtherServiceClient
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get note")
	defer span.Finish()

	span.SetTag("id", req.GetId())

	note, err := s.otherServiceClient.Get(ctx, 15)
	if err != nil {
		return nil, errors.WithMessage(err, "getting note")
	}

	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: note.ID,
			Info: &desc.NoteInfo{
				Title:   note.Info.Title,
				Content: note.Info.Content,
			},
			CreatedAt: timestamppb.New(note.CreatedAt),
			UpdatedAt: updatedAt,
		},
	}, nil
}

func main() {
	tracing.Init(serviceName)

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", grpcOtherPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatalf("failed to connect other service: %v", err)
	}
	otherServiceClient := other_service.New(descOther.NewOtherNoteV1Client(conn))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen port")
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.ServerTraicingInterceptor,
	))
	reflection.Register(s)

	desc.RegisterNoteV1Server(s, &server{
		otherServiceClient: otherServiceClient,
	})

	log.Println("grpc server is running")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to listen port")
	}
}
