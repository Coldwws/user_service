package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"user_service/internal/api/user"
	"user_service/internal/model"
	"user_service/internal/service"
	"user_service/internal/service/mocks"
	desc "user_service/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx         = context.Background()
		mc          = minimock.NewController(t)
		id          = gofakeit.Int64()
		firstName   = gofakeit.FirstName()
		lastName    = gofakeit.LastName()
		password    = gofakeit.Password(true, true, true, true, true, 6)
		phoneNumber = gofakeit.Phone()
		email       = gofakeit.Email()

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Info: &desc.UserInfo{
				FirstName:   firstName,
				LastName:    lastName,
				Password:    password,
				PhoneNumber: phoneNumber,
				Email:       email,
			},
		}
		info = &model.User{
			FirstName: firstName,
			LastName:  lastName,
			Password:  password,
			Phone:     phoneNumber,
			Email:     email,
		}
		res = &desc.CreateResponse{
			Id: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{ctx: ctx, req: req},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},

		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)

			NewId, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, NewId)
		})
	}
}
