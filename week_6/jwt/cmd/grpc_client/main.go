package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	descAccess "github.com/airats/micro_as_bigtech_course/week_6/jwt/pkg/access_v1"
)

const (
	grpcPort = 50061
)

var accessToken = flag.String("a", "", "access token")

func main() {
	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{
		"authorization": "Bearer" + *accessToken,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	client := descAccess.NewAccessV1Client(conn)

	_, err = client.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: "note_v1.NoteV1.Get",
	})
	if err != nil {
		log.Fatalf("failed to check: %v", err)
	}

	log.Println("Successful check")
}
