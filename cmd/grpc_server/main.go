package main

import (
	"context"
	"log"
	"net"
	"os"
	"user_service/internal/config"
	"user_service/internal/repository"
	"user_service/internal/repository/user"
	"user_service/internal/repository/user/model"
	desc "user_service/pkg/user_v1"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	userRepository repository.UserRepository
	desc.UnimplementedUserV1Server
}


func(s *server)Create(ctx context.Context, req *desc.CreateRequest)(*desc.CreateResponse,error){
	userModel := model.User{
		FirstName: req.Info.FirstName,
		LastName:  req.Info.LastName,
		Password:  req.Info.Password,
		Email:     req.Info.Email,
		Phone:     req.Info.PhoneNumber,
	}

	id,err := s.userRepository.Create(ctx, &userModel)
	if err!= nil{
		return nil,err
	}
	return &desc.CreateResponse{Id: id},nil

}


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
	srv := &server{
		userRepository: userRepo,
	}
	desc.RegisterUserV1Server(s, srv)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
