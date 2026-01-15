package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
	repo "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note"
	service "github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/note"
)

const (
	dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	noteRepo := repo.NewRepository(pool)

	serv := service.NewService(noteRepo)

	id, err := serv.Create(ctx, &model.NoteInfo{
		Title:   gofakeit.City(),
		Content: gofakeit.Animal(),
	})

	if err != nil {
		log.Fatalf("faild to created: %v", err)
	}

	fmt.Println(id)
}
