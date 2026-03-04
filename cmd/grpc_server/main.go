package main

import (
	"context"

	"log"
	"net"
	"os"
apiUser "user_service/internal/api/user"
	"user_service/internal/config"

	"user_service/internal/repository/user"
	userService "user_service/internal/service/user"

	desc "user_service/pkg/user_v1"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	if f := os.Getenv("ENV_FILE"); f != "" {
		_ = godotenv.Load(f)
	}
	cfg := config.LoadConfig()

	ctx := context.Background()

	poll, err := pgxpool.Connect(ctx, cfg.PG.DSN())
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer poll.Close()
	log.Println("Connected to Postgres")

	lis, err := net.Listen("tcp", cfg.GRPC.Addr())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	userRepo := user.NewRepository(poll)
	userSvc := userService.NewUserService(userRepo)
	
	srv := apiUser.NewServer(userSvc)


	desc.RegisterUserV1Server(s, srv)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
