package dto

type OrderRequest struct {
	ShopID          uint64               `json:"shop_id" validate:"required"` //customer
	ShippingAddress string               `json:"shipping_address" validate:"required"`
	Details         []OrderDetailRequest `json:"details" validate:"required"`
}

type OrderDetailRequest struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type PaymentRequest struct {
	OrderID       uint64  `json:"order_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}
