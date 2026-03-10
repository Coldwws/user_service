package user

import (
	"user_service/internal/service"
	desc "user_service/pkg/user_v1"
)
type Server struct{
	userService service.UserService
	desc.UnimplementedUserV1Server
}


func NewServer(userService service.UserService) *Server{
	return &Server{
		userService: userService,
	}
}
