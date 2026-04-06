package transaction

import (
	"context"
	"user_service/internal/client/db"
	"user_service/internal/client/db/pg"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct{
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager{
	return &manager{
		db: db,
	}
}

func(m *manager)transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler)(err error){

	tx,ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok{
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil{
		errors.Wrap(err,"Cannot Begin Transaction")
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func ()  {
		if r:= recover(); r != nil{
			err = errors.Errorf("panic recovered: %v",r)
		}

		if err != nil{
			if errRollBack := tx.Rollback(ctx); errRollBack != nil{
				err = errors.Wrap(errRollBack,"failed to rollback transaction")
			}
			return
		}

		if err == nil{
			err = tx.Commit(ctx)
			if err != nil{
				err = errors.Wrap(err,"failed to commit transaction")
			}
		}

	}()

	if err = fn(ctx); err != nil{
		err = errors.Wrap(err,"failed to executing code inside transaction")
	}
	return err

}

func (m *manager)ReadCommitted(ctx context.Context, f db.Handler)error{
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}