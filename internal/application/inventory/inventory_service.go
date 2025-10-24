package inventory

import (
	"context"
	"stocks/internal/domain"
	"stocks/internal/domain/models"
	"stocks/internal/shared/types"

	"github.com/beevik/guid"
)

type InventoryService interface {
	Inventory(ctx context.Context, domains []models.DeliveryItemer, path *types.Path) (*models.InventoryState, error)
}

type InventoryServiceImpl struct {
	uow domain.UnitOfWork
}

func NewInventoryService(uow domain.UnitOfWork) *InventoryServiceImpl {
	return &InventoryServiceImpl{uow: uow}
}

func (i *InventoryServiceImpl) Inventory(
	ctx context.Context,
	domains []models.DeliveryItemer,
	path *types.Path) (*models.InventoryState, error) {
	stockStates := make([]*models.StockState, 0)
	for _, d := range domains {
		states, err := d.Find(ctx, i.uow.Rest(), path)
		if err != nil {
			return nil, err
		}
		stockStates = append(stockStates, states...)
	}

	return withState(stockStates), nil
}

func withState(stockStates []*models.StockState) *models.InventoryState {
	var result models.InventoryResult
	nodes := make(map[guid.Guid]int)
	var toProduce int
	for _, stockState := range stockStates {
		if stockState.Produce {
			toProduce++
			continue
		}
		_, ok := nodes[*stockState.WarehouseID]
		if !ok {
			nodes[*stockState.WarehouseID] = 0
		}
	}
	if len(nodes) == 1 {
		result = models.AllInStockAtOne
	}
	if len(nodes) > 1 {
		result = models.AllInStockAtSeveral
	}
	if toProduce > 0 && len(nodes) > 0 {
		result = models.PartiallyInStock
	}
	if toProduce > 0 && len(nodes) == 0 {
		result = models.AllToProduce
	}

	return &models.InventoryState{Result: result, StockStates: stockStates}
}
