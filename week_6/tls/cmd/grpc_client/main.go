package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	desc "github.com/AiratS/micro_as_bigtech_course/week_1/grpc/pkg/note_v1"
	"github.com/fatih/color"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("service.pem", "")
	if err != nil {
		log.Fatalf("failed to load tls cer: %v", err)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.Get(ctx, &desc.GetRequest{Id: noteID})
	if err != nil {
		log.Fatalf("failed to get note: %v", err)
	}

	log.Println(color.RedString("Note info\n"), color.GreenString("%+v", resp.GetNote()))
}
