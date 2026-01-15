package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AiratS/micro_as_bigtech_course/week_3/config"
	impl "github.com/AiratS/micro_as_bigtech_course/week_3/internal/api/note"
	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	repo "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note"
	noteServ "github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/note"
)

const (
	grpcPort = 50061
	dbDSN    = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	config.Load(".env")

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

func newServer() *impl.Implementation {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	// defer pool.Close() TODO: close pool

	noteRepo := repo.NewRepository(pool)

	noteService := noteServ.NewService(noteRepo)

	impl.NewImplementation(noteService)

	return impl.NewImplementation(noteService)
}
