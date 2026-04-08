package main

import (
	"context"
	"google.golang.org/grpc/credentials"
	"log"
	"time"

	desc "user_service/pkg/user_v1"

	"github.com/fatih/color"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
	userID  = 13
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("Failed to connect server: #{err}")
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("Failed to get user by id: %v", err)
	}
	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", r.GetUser()))

	l, err := c.List(ctx, &desc.ListRequest{Limit: 12, Offset: 0})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	log.Printf(color.RedString("User list:\n"), color.GreenString("%+v", l.GetUsers()))
}
