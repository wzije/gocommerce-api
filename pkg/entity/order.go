package entity

import "time"

type OrderTable interface {
	TableName() string
}

func (Order) TableName() string {
	return "orders"
}

type Order struct {
	BaseEntity
	UserID          uint64        `gorm:"not null" json:"user_id" validate:"required"` //customer
	ShopID          uint64        `gorm:"not null" json:"shop_id" validate:"required"` //customer
	Date            time.Time     `gorm:"type:timestamp;default:current_timestamp" json:"date"`
	Status          string        `gorm:"type:varchar(50);not null" json:"status" validate:"required"`
	Amount          float64       `gorm:"type:decimal(10,2);not null" json:"amount" validate:"required"`
	ShippingDate    time.Time     `gorm:"type:timestamp" json:"shipping_date"`
	ShippingAddress string        `gorm:"type:text;not null" json:"shipping_address" validate:"required"`
	OrderDetail     []OrderDetail `gorm:"foreignkey:OrderID" json:"order_detail" validate:"required"`
}

type OrderDetailTable interface {
	TableName() string
}

func (OrderDetail) TableName() string {
	return "order_details"
}

type OrderDetail struct {
	BaseEntity
	OrderID      uint64   `gorm:"not null" json:"order_id" validate:"required"`
	ProductID    uint64   `gorm:"not null" json:"product_id" validate:"required"`
	Quantity     int      `gorm:"not null;check:quantity > 0" json:"quantity" validate:"required,gt=0"`
	PricePerUnit float64  `gorm:"type:decimal(10,2);not null" json:"price_per_unit" validate:"required"`
	TotalPrice   *float64 `gorm:"->;type:GENERATED ALWAYS AS (quantity*price_per_unit);" json:"total_price,omitempty"`
}
