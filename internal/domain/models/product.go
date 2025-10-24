package models

import (
	"context"
	"errors"

	"github.com/DimKa163/stock/internal/domain"
	"github.com/DimKa163/stock/internal/shared/collection"
	"github.com/DimKa163/stock/internal/shared/types"
	"github.com/beevik/guid"
	"github.com/shopspring/decimal"
)

type ChoicePriority int

const (
	Nearest ChoicePriority = iota
	Farthest
)

type DeliveryItemer interface {
	Find(ctx context.Context, restRepository domain.RestRepository, path *types.Path) ([]*StockState, error)
}

type InventoryProduct struct {
	ProductID    guid.Guid
	Quantity     decimal.Decimal
	IsLocal      bool
	IgnoredNodes []guid.Guid
}

type SimpleProduct struct {
	InventoryProduct
	ChoicePriority ChoicePriority
	FilialID       guid.Guid
}

func NewSimpleProduct(productID guid.Guid, quantity decimal.Decimal, isLocal bool, choice ChoicePriority, filialID guid.Guid) *SimpleProduct {
	return &SimpleProduct{
		InventoryProduct: InventoryProduct{
			ProductID: productID,
			Quantity:  quantity,
			IsLocal:   isLocal,
		},
		ChoicePriority: choice,
		FilialID:       filialID,
	}
}

func (sp *SimpleProduct) Find(ctx context.Context, restRepository domain.RestRepository, path *types.Path) ([]*StockState, error) {
	if path.Len() == 0 {
		return nil, errors.New("path is empty")
	}
	var strategy func(fn func(node guid.Guid) (bool, error)) error
	switch sp.ChoicePriority {
	case Nearest:
		strategy = path.Foreach
	case Farthest:
		strategy = path.ForeachReverse
	}
	remainingQuantity := sp.Quantity
	stockStates := make([]*StockState, 0)
	err := strategy(func(node guid.Guid) (bool, error) {
		rest, err := restRepository.Get(ctx, sp.FilialID, node, sp.ProductID)
		if err != nil {
			return false, err
		}
		if rest.Quantity.IsZero() {
			return true, nil
		}
		var covered decimal.Decimal
		if remainingQuantity.LessThanOrEqual(rest.Quantity) {
			covered = remainingQuantity
			remainingQuantity = remainingQuantity.Sub(covered)
		} else {
			covered = rest.Quantity
			remainingQuantity = remainingQuantity.Sub(covered)
		}
		stockStates = append(stockStates, &StockState{
			ProductID:   sp.ProductID,
			Quantity:    covered,
			WarehouseID: &node,
		})

		return remainingQuantity.GreaterThan(decimal.Zero), nil
	})

	if err != nil {
		return nil, err
	}

	if remainingQuantity.GreaterThan(decimal.Zero) {
		stockStates = append(stockStates, &StockState{
			ProductID: sp.ProductID,
			Quantity:  remainingQuantity,
			Produce:   true,
		})
	}
	return stockStates, nil
}

type CompositeProduct struct {
	Products       []*SimpleProduct
	ChoicePriority ChoicePriority
	FilialID       guid.Guid
}

func NewCompositeProduct(products []*SimpleProduct, choice ChoicePriority, filialID guid.Guid) *CompositeProduct {
	return &CompositeProduct{
		Products:       products,
		ChoicePriority: choice,
		FilialID:       filialID,
	}
}

func (cp *CompositeProduct) Find(ctx context.Context, restRepository domain.RestRepository, path *types.Path) ([]*StockState, error) {
	if path.Len() == 0 {
		return nil, errors.New("path is empty")
	}
	var strategy func(fn func(node guid.Guid) (bool, error)) error
	switch cp.ChoicePriority {
	case Nearest:
		strategy = path.Foreach
	case Farthest:
		strategy = path.ForeachReverse
	}
	stockStates := make([]*StockState, 0)
	allAtOne := false
	err := strategy(func(node guid.Guid) (bool, error) {
		remainingMap := make(map[guid.Guid]decimal.Decimal, len(cp.Products))
		for _, product := range cp.Products {
			remainingMap[product.ProductID] = product.Quantity
		}
		for _, product := range cp.Products {
			rest, err := restRepository.Get(ctx, cp.FilialID, node, product.ProductID)
			if err != nil {
				return false, err
			}
			if rest.Quantity.IsZero() {
				return true, nil
			}
			if remainingMap[product.ProductID].GreaterThan(rest.Quantity) {
				return true, nil
			}
			var covered decimal.Decimal
			if remainingMap[product.ProductID].LessThanOrEqual(rest.Quantity) {
				covered = remainingMap[product.ProductID]
			} else {
				covered = rest.Quantity
			}
			remainingMap[product.ProductID] = remainingMap[product.ProductID].Sub(covered)
		}
		covered := true
		for _, val := range collection.Values(remainingMap) {
			if !val.IsZero() {
				covered = false
			}
		}
		if covered {
			allAtOne = true
			for _, product := range cp.Products {
				stockStates = append(stockStates, &StockState{
					ProductID:   product.ProductID,
					Quantity:    product.Quantity,
					WarehouseID: &node,
				})
			}
		}
		return !covered, nil
	})
	if err != nil {
		return nil, err
	}
	if !allAtOne {
		for _, product := range cp.Products {
			product.ChoicePriority = Nearest
			if product.IsLocal {
				product.ChoicePriority = Farthest
			}
			stocks, err := product.Find(ctx, restRepository, path)
			if err != nil {
				return nil, err
			}
			for _, stock := range stocks {
				stockStates = append(stockStates, stock)
			}
		}
	}
	return stockStates, nil
}

type InventoryResult int

const (
	AllInStockAtOne InventoryResult = iota
	AllInStockAtSeveral
	PartiallyInStock
	AllToProduce
)

func (ir InventoryResult) String() string {
	return [...]string{"AllInStockAtOne", "AllInStockAtSeveral", "PartiallyInStock", "AllToProduce"}[ir]
}

type InventoryState struct {
	Result      InventoryResult
	StockStates []*StockState
}

type StockState struct {
	ProductID   guid.Guid
	Quantity    decimal.Decimal
	WarehouseID *guid.Guid
	Produce     bool
}
