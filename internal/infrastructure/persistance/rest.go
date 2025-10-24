package persistance

import (
	"context"
	"errors"
	"stocks/internal/domain"
	"stocks/internal/shared/db"

	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5"
)

const (
	restQuery = `SELECT id, quantity, filial_id, integration_ID, warehouse_id, product_id FROM public.rest
				WHERE filial_id = $1 AND warehouse_id = $2 AND product_id = $3`
)

type RestRepository struct {
	db db.QueryExecutor
}

func NewRestRepository(db db.QueryExecutor) RestRepository {
	return RestRepository{db: db}
}

func (r *RestRepository) Get(ctx context.Context, filialID guid.Guid, warehouseID guid.Guid, productID guid.Guid) (*domain.Rest, error) {
	var rest domain.Rest
	if err := r.db.QueryRow(ctx, restQuery, filialID, warehouseID, productID).Scan(&rest); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRestNotFound
		}
		return nil, err
	}
	return &rest, nil
}
