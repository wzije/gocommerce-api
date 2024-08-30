package entity

type ProductTable interface {
	TableName() string
}

func (Product) TableName() string { return "products" }

type Product struct {
	BaseEntity
	Name        string  `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gt=0"`
	SKU         string  `gorm:"type:varchar(100);unique;not null" json:"sku" validate:"required"`
	ImageURL    string  `gorm:"type:text" json:"image_url"`
	//Stock       *int    `json:"stock"`
	ShopID uint64 `json:"shop_id" validate:"required"`
}
