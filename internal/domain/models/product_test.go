package models

import (
	"context"
	"fmt"
	"testing"

	"github.com/DimKa163/stock/internal/domain"
	"github.com/DimKa163/stock/internal/shared/types"
	"github.com/DimKa163/stock/mocks"
	"github.com/beevik/guid"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSimpleProduct(t *testing.T) {
	var cases = []struct {
		Name    string
		Arrange func() (context.Context, *SimpleProduct, domain.RestRepository, *types.Path, []*StockState)
	}{
		{
			Name: "Simple Product should pick up from nearest warehouse",
			Arrange: func() (ctx context.Context, prd *SimpleProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)
				ctx = context.Background()
				prdID := guid.New()
				s := prdID.String()
				fmt.Println(s)
				filialID := guid.New()
				warehouseID := guid.New()
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, *filialID, *warehouseID, *prdID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(100),
					ProductID:     *prdID,
					WarehouseID:   *warehouseID,
				}, nil)
				rep = mockRep
				path := types.NewPath(1)
				path.AddNode(*warehouseID)
				p = path
				prd = NewSimpleProduct(*prdID, decimal.NewFromInt(3), true, Nearest, *filialID)
				exp = []*StockState{
					{
						WarehouseID: warehouseID,
						Quantity:    decimal.NewFromInt(3),
						ProductID:   *prdID,
					},
				}
				return
			},
		},
		{
			Name: "Simple Product should pick up from nearest warehouses",
			Arrange: func() (ctx context.Context, prd *SimpleProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)
				ctx = context.Background()
				prdID := guid.New()
				s := prdID.String()
				fmt.Println(s)
				filialID := guid.New()
				warehouseID1 := guid.New()
				warehouseID2 := guid.New()
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, *filialID, *warehouseID1, *prdID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(1),
					ProductID:     *prdID,
					WarehouseID:   *warehouseID1,
				}, nil)
				mockRep.EXPECT().Get(ctx, *filialID, *warehouseID2, *prdID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(2),
					ProductID:     *prdID,
					WarehouseID:   *warehouseID2,
				}, nil)
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(*warehouseID1)
				path.AddNode(*warehouseID2)
				p = path
				prd = NewSimpleProduct(*prdID, decimal.NewFromInt(3), false, Nearest, *filialID)
				exp = []*StockState{
					{
						WarehouseID: warehouseID1,
						Quantity:    decimal.NewFromInt(1),
						ProductID:   *prdID,
					},
					{
						WarehouseID: warehouseID2,
						Quantity:    decimal.NewFromInt(2),
						ProductID:   *prdID,
					},
				}
				return
			},
		},
		{
			Name: "Simple Product should pick up from nearest warehouses/with producing",
			Arrange: func() (ctx context.Context, prd *SimpleProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)
				ctx = context.Background()
				prdID := guid.New()
				s := prdID.String()
				fmt.Println(s)
				filialID := guid.New()
				warehouseID1 := guid.New()
				warehouseID2 := guid.New()
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, *filialID, *warehouseID1, *prdID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(1),
					ProductID:     *prdID,
					WarehouseID:   *warehouseID1,
				}, nil)
				mockRep.EXPECT().Get(ctx, *filialID, *warehouseID2, *prdID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(2),
					ProductID:     *prdID,
					WarehouseID:   *warehouseID2,
				}, nil)
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(*warehouseID1)
				path.AddNode(*warehouseID2)
				p = path
				prd = NewSimpleProduct(*prdID, decimal.NewFromInt(5), false, Nearest, *filialID)
				exp = []*StockState{
					{
						WarehouseID: warehouseID1,
						Quantity:    decimal.NewFromInt(1),
						ProductID:   *prdID,
					},
					{
						WarehouseID: warehouseID2,
						Quantity:    decimal.NewFromInt(2),
						ProductID:   *prdID,
					},
					{
						Quantity:  decimal.NewFromInt(2),
						ProductID: *prdID,
						Produce:   true,
					},
				}
				return
			},
		},
		{
			Name: "Simple Product should pick up from farthest warehouse",
			Arrange: func() (ctx context.Context, prd *SimpleProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)
				ctx = context.Background()
				prdID := guid.New()
				s := prdID.String()
				fmt.Println(s)
				filialID := guid.New()
				warehouseID := guid.New()
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, *filialID, *warehouseID, *prdID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(100),
					ProductID:     *prdID,
					WarehouseID:   *warehouseID,
				}, nil)
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(*guid.New())
				path.AddNode(*warehouseID)
				p = path
				prd = NewSimpleProduct(*prdID, decimal.NewFromInt(3), false, Farthest, *filialID)
				exp = []*StockState{
					{
						WarehouseID: warehouseID,
						Quantity:    decimal.NewFromInt(3),
						ProductID:   *prdID,
					},
				}
				return
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, prd, rep, p, exp := tt.Arrange()
			sut, err := prd.Find(ctx, rep, p)
			assert.NoError(t, err)
			assert.Equal(t, exp, sut)
		})
	}
}

func TestCompositeProduct(t *testing.T) {
	cases := []struct {
		Name    string
		Arrange func() (ctx context.Context, prd *CompositeProduct, rep domain.RestRepository, p *types.Path, exp []*StockState)
	}{
		{
			Name: "Composite Product should be all at one warehouse and pick up from nearest one",
			Arrange: func() (ctx context.Context, prd *CompositeProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)

				ctx = context.Background()
				filialID := *guid.New()
				warehouseID1 := *guid.New()
				warehouseID2 := *guid.New()
				prod1 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(10),
					},
					FilialID: filialID,
				}
				prod2 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(5),
					},
					FilialID: filialID,
				}
				prd = NewCompositeProduct([]*SimpleProduct{prod1, prod2}, Nearest, filialID)
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(100),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID1,
				}, nil)
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(50),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID1,
				}, nil)
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(warehouseID1)
				path.AddNode(warehouseID2)
				p = path
				exp = []*StockState{
					{
						ProductID:   prod1.ProductID,
						Quantity:    decimal.NewFromInt(10),
						WarehouseID: &warehouseID1,
					},
					{
						ProductID:   prod2.ProductID,
						Quantity:    decimal.NewFromInt(5),
						WarehouseID: &warehouseID1,
					},
				}
				return
			},
		},
		{
			Name: "Composite Product should be at one warehouse and pick up from nearest one but not first",
			Arrange: func() (ctx context.Context, prd *CompositeProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)

				ctx = context.Background()
				filialID := *guid.New()
				warehouseID1 := *guid.New()
				warehouseID2 := *guid.New()
				prod1 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(10),
					},
					FilialID: filialID,
				}
				prod2 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(5),
					},
					FilialID: filialID,
				}
				prd = NewCompositeProduct([]*SimpleProduct{prod1, prod2}, Nearest, filialID)
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID2, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(50),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID2, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(50),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(warehouseID1)
				path.AddNode(warehouseID2)
				p = path
				exp = []*StockState{
					{
						ProductID:   prod1.ProductID,
						Quantity:    decimal.NewFromInt(10),
						WarehouseID: &warehouseID2,
					},
					{
						ProductID:   prod2.ProductID,
						Quantity:    decimal.NewFromInt(5),
						WarehouseID: &warehouseID2,
					},
				}
				return
			},
		},
		{
			Name: "Composite Product should pick up from different warehouse",
			Arrange: func() (ctx context.Context, prd *CompositeProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)
				ctx = context.Background()
				filialID := *guid.New()
				warehouseID1 := *guid.New()
				warehouseID2 := *guid.New()
				prod1 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(10),
					},
					FilialID: filialID,
				}
				prod2 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(5),
					},
					FilialID: filialID,
				}
				prd = NewCompositeProduct([]*SimpleProduct{prod1, prod2}, Nearest, filialID)
				mockRep := mocks.NewMockRestRepository(ctrl)
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(100),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID2, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID2, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(20),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(warehouseID1)
				path.AddNode(warehouseID2)
				p = path
				exp = []*StockState{
					{
						ProductID:   prod1.ProductID,
						Quantity:    decimal.NewFromInt(10),
						WarehouseID: &warehouseID1,
					},
					{
						ProductID:   prod2.ProductID,
						Quantity:    decimal.NewFromInt(5),
						WarehouseID: &warehouseID2,
					},
				}
				return
			},
		},
		{
			Name: "Composite Product should pick up from different warehouse and have different direction",
			Arrange: func() (ctx context.Context, prd *CompositeProduct, rep domain.RestRepository, p *types.Path, exp []*StockState) {
				ctrl := gomock.NewController(t)
				ctx = context.Background()
				filialID := *guid.New()
				warehouseID1 := *guid.New()
				warehouseID2 := *guid.New()
				warehouseID3 := *guid.New()
				warehouseID4 := *guid.New()
				prod1 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(10),
					},
					FilialID: filialID,
				}
				prod2 := &SimpleProduct{
					InventoryProduct: InventoryProduct{
						ProductID: *guid.New(),
						Quantity:  decimal.NewFromInt(5),
						IsLocal:   true,
					},
					FilialID: filialID,
				}
				prd = NewCompositeProduct([]*SimpleProduct{prod1, prod2}, Nearest, filialID)
				mockRep := mocks.NewMockRestRepository(ctrl)
				// 1-st wh
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(100),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID1, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID1,
				}, nil).AnyTimes()
				// 2-nd wh
				mockRep.EXPECT().Get(ctx, filialID, warehouseID2, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID2, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(20),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				// 3-rd wh
				mockRep.EXPECT().Get(ctx, filialID, warehouseID3, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID3, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				// 4-rd wh
				mockRep.EXPECT().Get(ctx, filialID, warehouseID4, prod1.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(0),
					ProductID:     prod1.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				mockRep.EXPECT().Get(ctx, filialID, warehouseID4, prod2.ProductID).Return(&domain.Rest{
					RestID:        *guid.New(),
					FilialID:      &filialID,
					IntegrationID: guid.New(),
					Quantity:      decimal.NewFromInt(20),
					ProductID:     prod2.ProductID,
					WarehouseID:   warehouseID2,
				}, nil).AnyTimes()
				rep = mockRep
				path := types.NewPath(3)
				path.AddNode(warehouseID1)
				path.AddNode(warehouseID2)
				path.AddNode(warehouseID3)
				path.AddNode(warehouseID4)
				p = path
				exp = []*StockState{
					{
						ProductID:   prod1.ProductID,
						Quantity:    decimal.NewFromInt(10),
						WarehouseID: &warehouseID1,
					},
					{
						ProductID:   prod2.ProductID,
						Quantity:    decimal.NewFromInt(5),
						WarehouseID: &warehouseID4,
					},
				}
				return
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			ctx, prd, rep, p, exp := tt.Arrange()
			sut, err := prd.Find(ctx, rep, p)
			assert.NoError(t, err)
			assert.Equal(t, exp, sut)
		})
	}
}
