package main

import (
	"context"
	"log"
	"time"

	desc "github.com/AiratS/micro_as_bigtech_course/week_8/custom_errors/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:50061"

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect grpc: %v", err)
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: 10})
	if err != nil {
		log.Fatalf("here failed to get: %v", err)
	}

	log.Printf("%v", r.GetNote())
}
