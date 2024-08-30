package entity

type WarehouseTable interface {
	TableName() string
}

func (Warehouse) TableName() string {
	return "warehouses"
}

type Warehouse struct {
	BaseEntity
	Name               string               `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	IsActive           bool                 `gorm:"default:true" json:"is_active" validate:"required"`
	Location           string               `gorm:"type:varchar(255)" json:"location"`
	ShopID             uint64               `gorm:"type:varchar(255)" json:"shop_id"`
	UserID             uint64               `gorm:"type:varchar(255)" json:"user_id"`
	WarehouseInventory []WarehouseInventory `gorm:"foreignkey:WarehouseID" json:"warehouse_inventory"`
}

type WarehouseInventoryTable interface {
	TableName() string
}

func (WarehouseInventory) TableName() string {
	return "warehouse_inventories"
}

type WarehouseInventory struct {
	BaseEntity
	ProductID   uint64 `gorm:"not null" json:"product_id" validate:"required"`
	WarehouseID uint64 `gorm:"not null" json:"warehouse_id" validate:"required"`
	Quantity    int    `gorm:"not null;check:quantity >= 0" json:"quantity" validate:"required,gt=0"`
	//Products    *[]Product `gorm:"foreignkey:ProductID;references:ID" json:"products"`
	//Warehouse   *[]Warehouse `gorm:"foreignkey:WarehouseID;references:ID" json:"warehouses"`
}

type OrderWarehouseAllocationTable interface {
	TableName() string
}

func (OrderWarehouseAllocation) TableName() string {
	return "order_warehouse_allocations"
}

type OrderWarehouseAllocation struct {
	BaseEntity
	OrderID     uint64 `gorm:"not null" json:"order_id" validate:"required"`
	WarehouseID uint64 `gorm:"not null" json:"warehouse_id" validate:"required"`
	ProductID   uint64 `gorm:"not null" json:"product_id" validate:"required"`
	Quantity    int    `gorm:"not null" json:"quantity" validate:"required,min=1"`
}

type ProductTransferWarehouseTable interface {
	TableName() string
}

func (ProductTransferWarehouse) TableName() string {
	return "product_transfer_warehouses"
}

type ProductTransferWarehouse struct {
	BaseEntity
	SourceWarehouseID      uint64 `gorm:"not null" json:"source_warehouse_id"`
	DestinationWarehouseID uint64 `gorm:"not null" json:"destination_warehouse_id"`
	ProductID              uint64 `gorm:"not null" json:"product_id"`
	Quantity               int    `gorm:"not null" json:"quantity" validate:"gte=0"`
	Status                 string `gorm:"not null" json:"status"` // e.g., 'Pending', 'Completed', 'Failed'
}
