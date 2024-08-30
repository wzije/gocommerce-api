package entity

import "time"

type StockLockTable interface {
	TableName() string
}

func (StockLock) TableName() string {
	return "stock_locks"
}

type StockLock struct {
	BaseEntity
	OrderID     uint64     `gorm:"not null" json:"order_id" validate:"required"`
	WarehouseID uint64     `gorm:"not null" json:"warehouse_id" validate:"required"`
	ProductID   uint64     `gorm:"not null" json:"product_id" validate:"required"`
	Quantity    int        `gorm:"not null" json:"quantity" validate:"required,min=1"`
	Status      string     `gorm:"not null" json:"status" validate:"required"`
	LockedAt    time.Time  `gorm:"autoCreateTime" json:"locked_at"`
	ReleasedAt  *time.Time `gorm:"default:null" json:"released_at"`
}
