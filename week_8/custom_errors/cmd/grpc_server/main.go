package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys"
	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys/codes"
	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys/validate"
	"github.com/AiratS/micro_as_bigtech_course/week_8/custom_errors/internal/interceptor"
	desc "github.com/AiratS/micro_as_bigtech_course/week_8/custom_errors/pkg/note_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50061

type server struct {
	desc.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	time.Sleep(5 * time.Second)

	err := validate.Validate(
		ctx,
		validateID(req.GetId()),
		otherValidateID(req.GetId()),
	)
	if err != nil {
		return nil, err
	}

	if req.GetId() > 100 {
		return nil, sys.NewCommonError("id must be less than 100", codes.ResourceExhausted)
	}

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

func validateID(id int64) validate.Condition {
	return func(ctx context.Context) error {
		if id <= 0 {
			return validate.NewValidationErrors("id must be greater than 0")
		}

		return nil
	}
}

func otherValidateID(id int64) validate.Condition {
	return func(ctx context.Context) error {
		if id <= 100 {
			return validate.NewValidationErrors("id must be greater than 100")
		}

		return nil
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.ErrorCodesInterceptor,
	))
	desc.RegisterNoteV1Server(s, &server{})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
