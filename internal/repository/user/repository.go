package user

import (
	"context"
	"time"
	"user_service/internal/model"
	"user_service/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	Columns("first_name","last_name","phone_number","email","password","created_at").
	Values(user.FirstName,user.LastName,user.Phone,user.Email,user.Password,time.Now()).Suffix("RETURNING id")


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

func	(r *repo)	List(ctx context.Context, limit int64, offset int64)([]*model.User,error){
	qb := sq.Select("id","first_name","last_name","email","phone_number","created_at","updated_at").
	From("users").PlaceholderFormat(sq.Dollar).OrderBy("id").Limit(uint64(limit)).Offset(uint64(offset))

	query,args,err := qb.ToSql()
	if err != nil{
		return nil,err
	}

	rows,err := r.db.Query(ctx,query,args...)
	if err != nil{
		return nil,err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next(){
		var user model.User
			err = rows.Scan(
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
			users = append(users, &user)
		}
		return users,nil
}

func	(r *repo)	Update(ctx context.Context, id int64, updateUser *model.User)(error){
	qb := sq.Update("users").
        PlaceholderFormat(sq.Dollar).
        Where(sq.Eq{"id": id})

    
    if updateUser.FirstName != "" {
        qb = qb.Set("first_name", updateUser.FirstName)
    }
    if updateUser.LastName != "" {
        qb = qb.Set("last_name", updateUser.LastName)
    }
    if updateUser.Email != "" {
        qb = qb.Set("email", updateUser.Email)
    }
    if updateUser.Phone != "" {
        qb = qb.Set("phone_number", updateUser.Phone)
    }
    if updateUser.Password != "" {
        qb = qb.Set("password", updateUser.Password)
    }

    
    qb = qb.Set("updated_at", updateUser.UpdatedAt.Time)

    query, args, err := qb.ToSql()
    if err != nil {
        return err
    }

    res, err := r.db.Exec(ctx, query, args...)
    if err != nil {
        return err
    }

    if res.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }

    return nil
}


func	(r *repo)	Delete(ctx context.Context, id int64)(error){
	qb := sq.Delete("users").Where(sq.Eq{"id":id}).PlaceholderFormat(sq.Dollar)

	query,args,err := qb.ToSql()
	if err != nil{
		return err
	}

	del,err := r.db.Exec(ctx,query,args...)
	if err != nil{
		return err
	}

	if del.RowsAffected()==0{
		return pgx.ErrNoRows
	}

	return nil
}
