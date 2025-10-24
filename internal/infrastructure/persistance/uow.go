package persistance

import (
	"context"
	"stocks/internal/domain"
	"stocks/internal/shared/db"
)

type UnitOfWork struct {
	db db.QueryExecutor
}

func NewUnitOfWork(db db.QueryExecutor) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Rest() domain.RestRepository {
	return NewRestRepository(u.db)
}

func (u *UnitOfWork) Begin(ctx context.Context, fn func(work domain.UnitOfWork) error) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	err = fn(NewUnitOfWork(tx))

	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return nil
	}
	return tx.Commit(ctx)
}
