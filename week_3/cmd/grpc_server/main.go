package main

import (
	"context"
	"log"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/app"
)

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init the App: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
