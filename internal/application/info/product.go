package info

import (
	"github.com/beevik/guid"
	"github.com/shopspring/decimal"
)

type RequestedProduct struct {
	ProductID guid.Guid
	Quantity  decimal.Decimal
}
