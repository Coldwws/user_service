package user
import (
	"github.com/jackc/pgx/v4/pgxpool"
	"user_service/internal/repository"
	"context"
	sq "github.com/Masterminds/squirrel"
	"user_service/internal/repository/user/model"
)

type repo struct{
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository{
	return &repo{db: db}
}

func (r *repo)Create(ctx context.Context, user *model.User)(int64,error){
	qb := sq.Insert("users").
	PlaceholderFormat(sq.Dollar).
	Columns("first_name","last_name","phone_number","email","password").
	Values(user.FirstName,user.LastName,user.Phone,user.Email,user.Password).Suffix("RETURNING id")


	query,args,err := qb.ToSql()
	if err != nil { return 0,err}

	var userID int64

	err = r.db.QueryRow(ctx,query,args...).Scan(&userID)
	if err !=nil{return 0,err}

	return userID,nil

}

func (r *repo)	Get(ctx context.Context, id int64)(*model.User,error){
	qb := sq.Select("id","first_name","last_name","email","phone_number","created_at","updated_at").
	From("users").
	Where(sq.Eq{"id":id}).
	PlaceholderFormat(sq.Dollar)

	query,args,err := qb.ToSql()
	if err != nil{
		return nil,err
	}

	var user model.User

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil{
		return nil,err
	}

	return &user,nil
}

// func	(r *repo)	List(ctx context.Context, limit int64, offset int64)([]*model.User,error){}

// func	(r *repo)	Update(ctx context.Context, id int64, info *model.User)(error){}

// func	(r *repo)	Delete(ctx context.Context, id int64)(error){}
