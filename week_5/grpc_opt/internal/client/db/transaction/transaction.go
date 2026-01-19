package transaction

import (
	"context"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/db"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/pg"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type txManager struct {
	db db.Transactor
}

func NewTxManager(db db.Transactor) db.TxManager {
	return &txManager{
		db: db,
	}
}

func (t *txManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	opts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	return t.transaction(ctx, opts, f)
}

func (t *txManager) transaction(ctx context.Context, opts pgx.TxOptions, f db.Handler) (err error) {
	tx, ok := ctx.Value(db.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	tx, err = t.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	ctx = pg.MakeContextTx(ctx, tx)
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic after recovery: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "error after rollback: %v", errRollback)
			}

			return
		}

		if err == nil {
			errCommit := tx.Commit(ctx)
			if errCommit != nil {
				err = errors.Wrapf(err, "tx commit failed: %v", errCommit)
			}
		}
	}()

	if err = f(ctx); err != nil {
		return errors.Wrap(err, "can't call function inside transaction")
	}

	return nil
}
