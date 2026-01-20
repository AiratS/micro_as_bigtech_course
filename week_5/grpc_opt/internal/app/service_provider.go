package app

import (
	"context"
	"log"

	"github.com/AiratS/micro_as_bigtech_course/week_3/config"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/api/note"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/db"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/db/transaction"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/pg"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/closer"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository"
	noteRepo "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
	noteService "github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/note"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	dbClient       db.Client
	txManager      db.TxManager
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

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		httpConfig, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		sp.httpConfig = httpConfig
	}

	return sp.httpConfig
}

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		client, err := pg.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}

		closer.Add(client.Close)

		sp.dbClient = client
	}

	return sp.dbClient
}

func (sp *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if sp.noteRepository == nil {
		sp.noteRepository = noteRepo.NewRepository(sp.DBClient(ctx))
	}

	return sp.noteRepository
}

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTxManager(sp.DBClient(ctx).DB())
	}

	return sp.txManager
}

func (sp *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if sp.noteService == nil {
		sp.noteService = noteService.NewService(
			sp.NoteRepository(ctx),
			sp.TxManager(ctx),
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
