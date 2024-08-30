package entity

type ShopTable interface {
	TableName() string
}

func (Shop) TableName() string {
	return "shops"
}

type Shop struct {
	BaseEntity
	Name        string      `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Description string      `gorm:"type:text" json:"description"`
	Address     string      `gorm:"type:text" json:"address"`
	Phone       string      `gorm:"type:varchar(50)" json:"phone"`
	UserID      uint64      `gorm:"index" json:"user_id"`
	Warehouse   []Warehouse `gorm:"foreignkey:ShopID" json:"warehouse"`
	Products    []Product   `gorm:"foreignkey:ShopID" json:"products"`
}
