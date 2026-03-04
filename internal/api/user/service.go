package user

import (
	"user_service/internal/service"
	desc "user_service/pkg/user_v1"
)
type server struct{
	userService service.UserService
	desc.UnimplementedUserV1Server
}


func NewServer(userService service.UserService)*server{
	return &server{
		userService: userService,
	}
}
