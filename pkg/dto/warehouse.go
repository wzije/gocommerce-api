package dto

type WarehouseRequest struct {
	Name     string `json:"name" validate:"required"`
	ShopID   uint64 `json:"shop_id" validate:"required"`
	Location string `json:"location"`
}

type InventoryStockRequest struct {
	ProductID   uint64 `json:"product_id" validate:"required"`
	WarehouseID uint64 `json:"warehouse_id" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required"`
}

type TransferStockRequest struct {
	SourceWarehouseID uint64 `json:"source_warehouse_id" validate:"required"`
	TargetWarehouseID uint64 `json:"target_warehouse_id" validate:"required"`
	Quantity          int    `json:"quantity" validate:"required"`
	ProductID         uint64 `json:"product_id" validate:"required"`
}
