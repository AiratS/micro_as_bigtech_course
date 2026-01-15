package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/converter"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	repo "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note"
	impl "github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/note"
)

type server struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

func NewServer(noteService service.NoteService) *server {
	return &server{
		noteService: noteService,
	}
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.Info))
	if err != nil {
		return nil, err
	}

	log.Printf("Note created! %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	info, err := s.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("Get note data: %v", info)

	return &desc.GetResponse{
		Note: converter.ToDescNoteFromService(info),
	}, nil
}

const (
	grpcPort = 50061
	dbDSN    = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	s := newServer()
	desc.RegisterNoteV1Server(grpcServer, s)

	log.Printf("Server listening at: %v", lis.Addr())

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server GRPC server: %v", err)
	}
}

func newServer() *server {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	// defer pool.Close() TODO: close pool

	noteRepo := repo.NewRepository(pool)

	impl := impl.NewService(noteRepo)

	return &server{
		noteService: impl,
	}
}
