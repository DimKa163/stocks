package domain

import (
	"context"

	"github.com/beevik/guid"
	"github.com/shopspring/decimal"
)

type Rest struct {
	RestID        guid.Guid
	FilialID      *guid.Guid
	IntegrationID *guid.Guid
	Quantity      decimal.Decimal
	ProductID     guid.Guid
	WarehouseID   guid.Guid
}

type RestRepository interface {
	Get(ctx context.Context, filialID guid.Guid, warehouseID guid.Guid, productID guid.Guid) (*Rest, error)
}
