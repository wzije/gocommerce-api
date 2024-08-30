package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	if err := r.db.WithContext(ctx).Create(payment); err != nil {
		return err.Error
	}
	return nil
}

func (r *paymentRepository) GetPaymentByOrderId(ctx context.Context, orderId uint64) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.WithContext(ctx).First(&payment, "order_id=?", orderId).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, id uint64, status string) error {
	payment := entity.Payment{}
	if err := r.db.WithContext(ctx).First(&payment, id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func NewPaymentRepository(db *gorm.DB) repository.PaymentRepositoryInterface {
	return &paymentRepository{
		db: db,
	}
}
