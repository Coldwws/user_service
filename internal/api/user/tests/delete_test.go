package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"user_service/internal/api/user"
	"user_service/internal/service"
	"user_service/internal/service/mocks"
	desc "user_service/pkg/user_v1"
)

func TestDelete(t *testing.T) {
	type userServiceMock func(mc *minimock.Controller) service.UserService

	ctx := context.Background()
	mc := minimock.NewController(t)
	defer mc.Finish()

	id := gofakeit.Int64()
	serviceErr := fmt.Errorf("service error")

	req := &desc.DeleteRequest{Id: id}

	tests := []struct {
		name      string
		mockSetup userServiceMock
		wantErr   error
	}{
		{
			name: "success delete",
			mockSetup: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
			wantErr: nil,
		},
		{
			name: "service error",
			mockSetup: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
			wantErr: serviceErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			srv := user.NewServer(tt.mockSetup(mc))
			resp, err := srv.Delete(ctx, req)

			require.Equal(t, tt.wantErr, err)

			if err == nil {
				require.NotNil(t, resp)
			}
		})
	}
}
