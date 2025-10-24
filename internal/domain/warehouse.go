package domain

import "github.com/beevik/guid"

type Warehouse struct {
	WarehouseID     guid.Guid
	Type            WarehouseType
	DescriptorGroup string
	RestAvailable   bool
	PickupOnly      bool
}

type WarehouseType int

const (
	WarehouseTypeFREE WarehouseType = iota
	WarehouseTypeMAIN
	WarehouseTypeMAINWH
	WarehouseTypeSHOPPINGCENTER
)

func (wt WarehouseType) String() string {
	return [...]string{"FREE", "MAIN", "MAIN_WAREHOUSE", "SHOPPING_CENTER"}[wt]
}
