package domain

import (
	"context"
)

type UnitOfWork interface {
	Rest() RestRepository

	Begin(ctx context.Context, fn func(work UnitOfWork) error) error
}
