package user

import (
	
	"context"
	"user_service/internal/repository"
	"user_service/internal/model"

)

type uService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *uService {
	return &uService{userRepository: userRepository}
}

func(s *uService)Create(ctx context.Context, user *model.User)(int64,error){
	id,err := s.userRepository.Create(ctx,user)
	if err !=nil{
		return 0,err
	}
	return id,nil
}

func (s *uService)Get(ctx context.Context, id int64)(*model.User,error){
	user,err := s.userRepository.Get(ctx,id)
	if err != nil{
		return nil,err
	}
	return user,nil

}

func (s *uService)List(ctx context.Context, limit int64, offset int64)([]*model.User,error){
	users,err := s.userRepository.List(ctx,limit,offset)
	if err != nil{
		 return nil,err
	}
	return users,nil
}

func (s *uService)Update(ctx context.Context, id int64, user *model.User)(error){
	err := s.userRepository.Update(ctx,id,user)
	if err != nil{
		return err
	}
	return nil
}

func (s *uService)Delete(ctx context.Context, id int64)error{
	err := s.userRepository.Delete(ctx,id)
	if err != nil{
		return err
	}
	return nil


}
