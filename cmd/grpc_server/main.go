package main

import (
	"context"
	
	"log"
	"net"

	"sync"

	"user_service/internal/config"
	desc "user_service/pkg/user_v1"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	users   = make(map[int64]*desc.User)
	usersMu sync.Mutex
)



type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server)Delete(ctx context.Context, req *desc.DeleteRequest)(*emptypb.Empty, error){
	log.Printf("Delete User")

	usersMu.Lock()
	defer usersMu.Unlock()

	if _, ok := users[req.GetId()]; !ok{
		return nil, status.Error(codes.NotFound, "user not found")
	}

	delete(users,req.Id)
	return &emptypb.Empty{}, nil
}


func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Update User")



	if req.Info == nil{
		return nil, status.Error(codes.InvalidArgument,"update info is required")
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	user,ok := users[req.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	
	info := req.Info

	if info.GetFirstName() != nil {
		user.Info.FirstName = info.GetFirstName().GetValue()
	}
	if info.GetLastName() != nil {
		user.Info.LastName = info.GetLastName().GetValue()
	}
	if info.GetPassword() != nil {
		user.Info.Password = info.GetPassword().GetValue()
	}
	if info.GetPhoneNumber() != nil {
		user.Info.PhoneNumber = info.GetPhoneNumber().GetValue()
	}
	if info.GetEmail() != nil {
		user.Info.Email = info.GetEmail().GetValue()
	}
	user.UpdatedAt = timestamppb.Now()

	return &emptypb.Empty{}, nil

}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create User")

	info := req.GetInfo()
	
	if info == nil {
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}

	if info.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}
	if info.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	userID := int64(gofakeit.Uint32())
	now := timestamppb.Now()

	user := &desc.User{
		Id:        userID,
		Info:      info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	users[userID] = user

	return &desc.CreateResponse{
		Id: userID,
	}, nil

}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	usersMu.Lock()
	defer usersMu.Unlock()

	user, ok := users[req.GetId()]

	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &desc.GetResponse{
		User: user,
	}, nil

}

func (s *server) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	log.Printf("List users")

	usersMu.Lock()
	defer usersMu.Unlock()

	limit := int(req.GetLimit())
	offset := int(req.GetOffset())

	var list []*desc.User
	i := 0

	
	for _, user := range users {
		if i >= offset && len(list) < limit {
			list = append(list, user)
		}
		i++
	}

	return &desc.ListResponse{
		Users: list,
	}, nil
}


func main() {
	_ = godotenv.Load("local.env")

	cfg := config.LoadConfig()
	lis, err := net.Listen("tcp", cfg.GRPC.Addr())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
