package repository

import (
	"context"
	"user_service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User)(int64,error)
	Get(ctx context.Context, id int64)(*model.User,error)
	List(ctx context.Context, limit int64, offset int64)([]*model.User,error)
	Update(ctx context.Context, id int64, user *model.User)(error)
	Delete(ctx context.Context, id int64)(error)
}