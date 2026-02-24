package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"user_service/internal/config"
	desc "user_service/pkg/user_v1"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type server struct {
	db *pgxpool.Pool
	desc.UnimplementedUserV1Server
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

	qb := sq.Insert("users").
	PlaceholderFormat(sq.Dollar).
	Columns("first_name","last_name","phone_number","email","password").
	Values(info.GetFirstName(),info.GetLastName(),info.GetPhoneNumber(),info.GetEmail(),info.GetPassword()).Suffix("RETURNING id")

	query,args,err := qb.ToSql()
	if err != nil{
		log.Fatalf("Failed to build query: %v",err)
	}
	var userID int
	err = s.db.QueryRow(ctx,query,args...).Scan(&userID)
	if err != nil{
		log.Fatalf("Failed to execute query: %v", err)
	}

	log.Printf("Created user with ID: %d", userID)
	return &desc.CreateResponse{Id: int64(userID)},nil
}

func(s *server) Update(ctx context.Context, req *desc.UpdateRequest)(*emptypb.Empty,error){
	log.Printf("Update User")

	info := req.GetInfo()

	if info == nil{
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}
	qb := sq.Update("users").
	PlaceholderFormat(sq.Dollar).
	Where(sq.Eq{"id":req.GetId()})

	changed := false

	if v := info.GetFirstName(); v != nil{
		qb = qb.Set("first_name",v.GetValue());changed=true
	}
	if v := info.GetLastName(); v!= nil{
		qb = qb.Set("last_name",v.GetValue()); changed=true
	}
	if v:= info.GetPhoneNumber(); v!= nil{
		qb = qb.Set("phone_number",v.GetValue());changed=true
	}
	if v := info.GetEmail(); v!= nil{
		qb = qb.Set("email",v.GetValue()); changed=true
	}
	if v := info.GetPassword(); v!=nil{
		qb =qb.Set("password", v.GetValue());changed= true
	}

	if !changed{
		return &emptypb.Empty{},nil
	}

	qb = qb.Set("updated_at",sq.Expr("now()"))

	query,args,err := qb.ToSql()
	if err != nil{
		return nil, status.Error(codes.Internal,"Failed to execute query")
	}
	ct,err := s.db.Exec(ctx,query,args...)
	if err != nil{
		return nil, status.Error(codes.Internal,"Failed to execute query")
	}
	if ct.RowsAffected() == 0{
		return nil, status.Error(codes.NotFound,"User not found")
	}

	log.Printf("User with id:%d updated",req.GetId())
	return &emptypb.Empty{},nil
}

func(s *server) Delete(ctx context.Context,req *desc.DeleteRequest)(*emptypb.Empty,error){
	log.Printf("Delete User")

	qb := sq.Delete("users").
	PlaceholderFormat(sq.Dollar).
	Where(sq.Eq{"id":req.GetId()})

	query,args,err := qb.ToSql()
	if err != nil{
		return nil, status.Error(codes.Internal,"Failed to build query")

	}

	ct,err := s.db.Exec(ctx,query,args...)
	if err != nil{
		return nil, status.Error(codes.Internal,"Failed to execute query")
	}
	if ct.RowsAffected() == 0{
		return nil, status.Error(codes.NotFound,"User not found")
	}
	log.Printf("User with id:%d deleted",req.GetId())
	return &emptypb.Empty{},nil
}

func (s *server)List(ctx context.Context, req *desc.ListRequest)(*desc.ListResponse,error){
	log.Printf("Users List")
	limit := uint64(req.GetLimit())
	if limit == 0{
		limit = 10
	}
	offset := uint64(req.GetOffset())

	qb := sq.Select("id","first_name","last_name","email","phone_number","created_at","updated_at").
	From("users").
	PlaceholderFormat(sq.Dollar).
	OrderBy("id ASC").
	Limit(limit).
	Offset(offset)

	query,args,err := qb.ToSql()
	if err != nil{
		log.Fatalf("Failed to build query: %v",err)
	}

	rows,err := s.db.Query(ctx,query,args...)
	if err != nil{
		log.Fatalf("Failed to execute query :%v",err)
	}
	defer rows.Close()
	users := make([]*desc.User,0,limit)
	for rows.Next(){
		u := &desc.User{Info: &desc.UserInfo{}}
		var createdAt,updatedAt sql.NullTime

		err := rows.Scan(&u.Id,&u.Info.FirstName,&u.Info.LastName,&u.Info.Email,&u.Info.PhoneNumber,&createdAt,&updatedAt)
		if err != nil{
			return nil,err
		}
		if createdAt.Valid{
			u.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			u.UpdatedAt = timestamppb.New(updatedAt.Time)
		}
		users = append(users,u)
	}

	return &desc.ListResponse{Users: users},nil
} 

func (s *server) Get(ctx context.Context, req *desc.GetRequest)(*desc.GetResponse,error){
	log.Printf("User id: %d", req.GetId())

	qb := sq.Select("id","first_name","last_name","email","phone_number","created_at","updated_at").
	From("users").
	Where(sq.Eq{"id":req.GetId()}).
	PlaceholderFormat(sq.Dollar)

	query,args,err := qb.ToSql()
	if err != nil{
		log.Fatalf("Failed to build query: %v",err)
	}
	user := desc.User{
		Info: &desc.UserInfo{
		},
	}
	var createdAt,updatedAt sql.NullTime

	err = s.db.QueryRow(ctx,query,args...).Scan(&user.Id,&user.Info.FirstName,&user.Info.LastName,&user.Info.Email,&user.Info.PhoneNumber,&createdAt,&updatedAt)
	if err != nil{
		log.Fatalf("Failed to execute query:%v",err)
	}

	if createdAt.Valid{
		user.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid{
		user.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	return &desc.GetResponse{User: &user},nil
}


func main() {
	if f:= os.Getenv("ENV_FILE"); f!= ""{
		_ = godotenv.Load(f)
	}

	cfg := config.LoadConfig()

	ctx := context.Background()

	poll,err := pgxpool.Connect(ctx,cfg.PG.DSN())
	if err != nil{
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

	srv := &server{db:poll}
	desc.RegisterUserV1Server(s, srv)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
