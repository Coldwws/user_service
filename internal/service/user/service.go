package user

import (
	
	"context"
	"user_service/internal/repository"
	"user_service/internal/model"
	"user_service/internal/client/db"

)

type uService struct {
	userRepository repository.UserRepository
	txManager db.TxManager
}

func NewUserService(userRepository repository.UserRepository, txManager db.TxManager) *uService {
	return &uService{userRepository: userRepository, txManager: txManager}
}

func(s *uService)Create(ctx context.Context, user *model.User)(int64,error){

	
	var id int64
	err := s.txManager.ReadCommitted(ctx, func (ctx context.Context)error  {
		var errTx error
		id, errTx = s.userRepository.Create(ctx,user)
			return errTx
	})
	return id,err
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
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context)error{
		return s.userRepository.Update(ctx,id,user)
	})

}

func (s *uService)Delete(ctx context.Context, id int64)error{
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context)error{
		return s.userRepository.Delete(ctx,id)
	})
}
