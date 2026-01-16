package app

import (
	"context"
	"log"

	"github.com/AiratS/micro_as_bigtech_course/week_3/config"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/api/note"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/closer"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository"
	noteRepo "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
	noteService "github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/note"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbPool         *pgxpool.Pool
	noteRepository repository.NoteRepository

	noteService service.NoteService

	noteImpl *note.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) PGConfig() config.PGConfig {
	if sp.pgConfig == nil {
		pgConfig, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("Failed to get PGConfig: %v", err)
		}

		sp.pgConfig = pgConfig
	}

	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		gprcConfig, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("Failed to get GRPCConfig: %v", err)
		}

		sp.grpcConfig = gprcConfig
	}

	return sp.grpcConfig
}

func (sp *serviceProvider) DBPool(ctx context.Context) *pgxpool.Pool {
	if sp.dbPool == nil {
		pool, err := pgxpool.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		sp.dbPool = pool
	}

	return sp.dbPool
}

func (sp *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if sp.noteRepository == nil {
		sp.noteRepository = noteRepo.NewRepository(sp.DBPool(ctx))
	}

	return sp.noteRepository
}

func (sp *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if sp.noteService == nil {
		sp.noteService = noteService.NewService(
			sp.NoteRepository(ctx),
		)
	}

	return sp.noteService
}

func (sp *serviceProvider) NoteImpl(ctx context.Context) *note.Implementation {
	if sp.noteImpl == nil {
		sp.noteImpl = note.NewImplementation(
			sp.NoteService(ctx),
		)
	}

	return sp.noteImpl
}
