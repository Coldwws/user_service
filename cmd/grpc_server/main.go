package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"
	"user_service/internal/config"
	"user_service/internal/converter"
	"user_service/internal/model"
	"user_service/internal/repository"
	"user_service/internal/repository/user"
	desc "user_service/pkg/user_v1"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	userRepository repository.UserRepository
	desc.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userModel := converter.UserProtoToModel(&desc.User{
		Info: &desc.UserInfo{
			FirstName:   req.Info.FirstName,
			LastName:    req.Info.LastName,
			PhoneNumber: req.Info.PhoneNumber,
			Email:       req.Info.Email,
		},
	})
	userModel.Password = req.Info.Password

	id, err := s.userRepository.Create(ctx, userModel)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil

}
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userModel, err := s.userRepository.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.UserModelToProto(userModel),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	del := s.userRepository.Delete(ctx, req.Id)
	if del != nil {
		return nil, del
	}
	return &emptypb.Empty{}, nil
}

func (s *server)List(ctx context.Context, req *desc.ListRequest)(*desc.ListResponse,error){
	userList, err := s.userRepository.List(ctx,req.Limit,req.Offset)
	if err != nil{
		 return &desc.ListResponse{}, err
	}
	protoUsers := make([]*desc.User,0,len(userList))
	for _,u := range userList{
		protoUsers = append(protoUsers,converter.UserModelToProto(u))
	}

	return &desc.ListResponse{
		Users: protoUsers,
	},nil
}

func (s *server)Update(ctx context.Context, req *desc.UpdateRequest)(*emptypb.Empty,error){
    userModel := &model.User{
        ID:        req.Id,
				UpdatedAt: sql.NullTime{Time: time.Now(),Valid: true},
    }

		    if req.Info.FirstName != nil {
        userModel.FirstName = req.Info.FirstName.GetValue()
    }
    if req.Info.LastName != nil {
        userModel.LastName = req.Info.LastName.GetValue()
    }
    if req.Info.Email != nil {
        userModel.Email = req.Info.Email.GetValue()
    }
    if req.Info.PhoneNumber != nil {
        userModel.Phone = req.Info.PhoneNumber.GetValue()
    }
    if req.Info.Password != nil {
        userModel.Password = req.Info.Password.GetValue()
    }
		
    err := s.userRepository.Update(ctx, req.Id, userModel)
    if err != nil {
        return nil, err
    }

    return &emptypb.Empty{}, nil
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
