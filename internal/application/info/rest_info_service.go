package info

import (
	"context"

	"github.com/DimKa163/stocks/internal/domain"
	"github.com/beevik/guid"
)

type RestInfoService interface {
	GetStockOneItemInfo(ctx context.Context, product RequestedProduct, filialID guid.Guid, shipment guid.Guid) (*OneStockInfo, error)

	GetStockManyItemsInfo(ctx context.Context, products []RequestedProduct, filialID, shipment guid.Guid) (*ManyStockInfo, error)
}

type RestInfoServiceImpl struct {
	uow domain.UnitOfWork
}

func NewRestInfoService(uow domain.UnitOfWork) *RestInfoServiceImpl {
	return &RestInfoServiceImpl{uow: uow}
}

func (r *RestInfoServiceImpl) GetStockOneItemInfo(ctx context.Context, product RequestedProduct, filialID, shipmentID guid.Guid) (*OneStockInfo, error) {
	restRepository := r.uow.Rest()
	rest, err := restRepository.Get(ctx, product.ProductID, filialID, shipmentID)
	if err != nil {
		return nil, err
	}
	if rest.Quantity.LessThan(product.Quantity) {
		return &OneStockInfo{
			InStock:     false,
			ProductInfo: ProductInfo{ProductID: product.ProductID},
		}, nil
	}
	return &OneStockInfo{
		InStock:     true,
		ProductInfo: ProductInfo{ProductID: product.ProductID, Rest: rest, Covered: true},
	}, nil
}

func (r *RestInfoServiceImpl) GetStockManyItemsInfo(ctx context.Context, products []RequestedProduct, filialID, shipment guid.Guid) (*ManyStockInfo, error) {
	description := make([]ProductInfo, len(products))
	allInStock := true
	for i, product := range products {
		stockInfo, err := r.GetStockOneItemInfo(ctx, product, filialID, shipment)
		if err != nil {
			return nil, err
		}
		if !stockInfo.InStock {
			allInStock = false
		}
		description[i] = stockInfo.ProductInfo
	}
	return &ManyStockInfo{
		InStock:     allInStock,
		ProductInfo: description,
	}, nil
}
