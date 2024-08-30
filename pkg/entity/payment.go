package entity

type PaymentTable interface {
	TableName() string
}

func (Payment) TableName() string {
	return "payments"
}

type Payment struct {
	BaseEntity
	OrderID       uint64  `gorm:"not null" json:"order_id" validate:"required"`
	PaymentMethod string  `gorm:"not null" json:"payment_method" validate:"required"`
	Amount        float64 `gorm:"type:decimal(10,2);not null" json:"amount" validate:"required"`
	Status        string  `gorm:"type:varchar(50);not null" json:"status" validate:"required,oneof=PENDING PAID FAILED"`
}
