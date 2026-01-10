package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
)

const (
	dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer con.Close(ctx)

	// Insert
	res, err := con.Exec(ctx, "INSERT INTO note(title, body) VALUES($1, $2)", gofakeit.City(), gofakeit.Address().Street)
	if err != nil {
		log.Fatalf("failed to insert to database: %v", err)
	}

	log.Printf("inserted rows: %v", res.RowsAffected())

	// Fetch
	rows, err := con.Query(ctx, "SELECT id, title, body, created_at, updated_at FROM note")
	if err != nil {
		log.Fatalf("failed to query rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		log.Printf("id: %d, title: %s, body: %s, createdAt: %v, updatedAt: %v", id, title, body, createdAt, updatedAt)
	}
}
