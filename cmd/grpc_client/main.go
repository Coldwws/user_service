package main

import (
	"context"
	"log"
	"time"

	desc "user_service/pkg/user_v1"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	userID = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil{
		log.Fatalf("Failed to connect server: #{err}")
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)
	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	r,err := c.Get(ctx,&desc.GetRequest{Id: userID})
	if err != nil{
		log.Fatalf("Failed to get user by id: %v",err)
	}
	log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", r.GetUser()))

	l,err := c.List(ctx,&desc.ListRequest{Limit: 5})
	if err != nil{
		log.Fatalf("Failed to list users: %v",err)
	}
	log.Printf(color.RedString("User list:\n"), color.GreenString("%+v", l.GetUsers()))
	}
