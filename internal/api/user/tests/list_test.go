package tests

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"user_service/internal/api/user"
	"user_service/internal/converter"
	"user_service/internal/model"
	"user_service/internal/service"
	"user_service/internal/service/mocks"
	desc "user_service/pkg/user_v1"
)

func TestList(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.ListRequest
	}

	ctx := context.Background()
	mc := minimock.NewController(t)

	defer mc.Finish()

	limit := int64(2)
	offset := int64(0)
	serviceErr := errors.New("service error")

	now := time.Now()
	user1 := &model.User{
		ID:        gofakeit.Int64(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Phone:     gofakeit.Phone(),
		Password:  gofakeit.Password(true, true, true, true, true, 6),
		CreatedAt: now,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}
	user2 := &model.User{
		ID:        gofakeit.Int64(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Phone:     gofakeit.Phone(),
		Password:  gofakeit.Password(true, true, true, true, true, 6),
		CreatedAt: now,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	req := &desc.ListRequest{Limit: limit, Offset: offset}

	tests := []struct {
		name            string
		args            args
		want            *desc.ListResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{ctx: ctx, req: req},
			want: &desc.ListResponse{
				Users: []*desc.User{
					converter.UserModelToProto(user1),
					converter.UserModelToProto(user2),
				},
			},
			err: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.ListMock.Expect(ctx, limit, offset).Return([]*model.User{user1, user2}, nil)
				return mock
			},
		},
		{
			name: "empty list",
			args: args{ctx: ctx, req: req},
			want: &desc.ListResponse{
				Users: []*desc.User{},
			},
			err: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.ListMock.Expect(ctx, limit, offset).Return([]*model.User{}, nil)
				return mock
			},
		},
		{
			name: "service error",
			args: args{ctx: ctx, req: req},
			want: &desc.ListResponse{
				Users: []*desc.User{},
			},
			err: serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.ListMock.Expect(ctx, limit, offset).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)
			resp, err := api.List(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.Users, resp.Users)
			}
		})
	}
}
