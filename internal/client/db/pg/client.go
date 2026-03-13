package pg

import (
	"context"
	"user_service/internal/client/db"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct{
	masterDBC db.DB
}

func New(ctx context.Context, dsn string) (db.Client,error){
	dbc, err := pgxpool.Connect(ctx,dsn)
	if err != nil{
		return nil, errors.Errorf("failed to connect to db: %v",err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	},nil
}

func(c *pgClient)DB() db.DB{
	return c.masterDBC
}

func MakeContextTx(ctx context.Context, tx pgx.Tx)context.Context{
	return context.WithValue(ctx,TxKey,tx)
}



func (c *pgClient)Close()error{
	if c.masterDBC != nil{
		c.masterDBC.Close()
	}

	return nil
}
