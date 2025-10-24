package info

import (
	"github.com/DimKa163/stock/internal/domain"
	"github.com/beevik/guid"
)

type OneStockInfo struct {
	InStock     bool
	ProductInfo ProductInfo
}

type ManyStockInfo struct {
	InStock     bool
	ProductInfo []ProductInfo
}

type ProductInfo struct {
	ProductID guid.Guid
	Rest      *domain.Rest
	Covered   bool
}
