package dto

type ProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	SKU         string  `json:"sku"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
}

type ProductResponse struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	ShopID      uint64  `json:"shop_id"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
}
