package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
	"user_service/internal/api/user"
	"user_service/internal/model"
	"user_service/internal/service"
	"user_service/internal/service/mocks"
	desc "user_service/pkg/user_v1"
)

func TestGet(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
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

		serviceErr      = fmt.Errorf("service error")
		userNotFoundErr = errors.New("user not found")
		req             = &desc.GetRequest{
			Id: id,
		}

		userModel = &model.User{
			ID:        id,
			FirstName: firstName,
			LastName:  lastName,
			Password:  password,
			Phone:     phoneNumber,
			Email:     email,
		}

		res = &desc.GetResponse{
			User: &desc.User{
				Id: id,
				Info: &desc.UserInfo{
					FirstName:   firstName,
					LastName:    lastName,
					Email:       email,
					PhoneNumber: phoneNumber,
				},
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(userModel, nil)
				return mock
			},
		},
		//Test #3
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
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "user not found",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  userNotFoundErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, userNotFoundErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)
			resp, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			if resp != nil && tt.want != nil {
				require.Equal(t, tt.want.User.Id, resp.User.Id)
				require.Equal(t, tt.want.User.Info.FirstName, resp.User.Info.FirstName)
				require.Equal(t, tt.want.User.Info.LastName, resp.User.Info.LastName)
				require.Equal(t, tt.want.User.Info.Email, resp.User.Info.Email)
				require.Equal(t, tt.want.User.Info.PhoneNumber, resp.User.Info.PhoneNumber)
			}
		})
	}
}
