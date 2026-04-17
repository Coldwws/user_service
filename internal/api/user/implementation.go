package user

import (
	"context"
	"database/sql"
	"math/rand"
	"time"
	"user_service/internal/converter"
	"user_service/internal/logger"
	"user_service/internal/model"
	desc "user_service/pkg/user_v1"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Create user", zap.Any("UserInfo", req.GetInfo()))

	userModel := converter.UserProtoToModel(&desc.User{
		Info: &desc.UserInfo{
			FirstName:   req.Info.FirstName,
			LastName:    req.Info.LastName,
			PhoneNumber: req.Info.PhoneNumber,
			Email:       req.Info.Email,
		},
	})

	userModel.Password = req.Info.Password

	id, err := s.userService.Create(ctx, userModel)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}
func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Get User", zap.Int64("UserID", req.Id))

	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetUser")
	defer span.Finish()

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	userModel, err := s.userService.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.UserModelToProto(userModel),
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	del := s.userService.Delete(ctx, req.Id)
	if del != nil {
		return nil, del
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetUserList")
	defer span.Finish()

	userList, err := s.userService.List(ctx, req.Limit, req.Offset)

	if err != nil {
		return &desc.ListResponse{}, err
	}

	protoUsers := make([]*desc.User, 0, len(userList))
	for _, u := range userList {
		protoUsers = append(protoUsers, converter.UserModelToProto(u))
	}

	return &desc.ListResponse{
		Users: protoUsers,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	userModel := &model.User{
		ID:        req.Id,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
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

	err := s.userService.Update(ctx, req.Id, userModel)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
