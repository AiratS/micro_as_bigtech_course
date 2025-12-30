package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/AiratS/micro_as_bigtech_course/week_1/grpc/pkg/note_v1"
	"github.com/fatih/color"
)

const (
	address = "localhost:50081"
	noteID  = 12
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
