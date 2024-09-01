package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func (r *orderRepository) MyOrders(ctx context.Context) (*[]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Where("orders.user_id", security.PayloadData.UserID).
		Find(&orders).Error
	return &orders, err
}

func (r *orderRepository) MyCustomerOrders(ctx context.Context) (*[]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Joins("LEFT JOIN shops ON shops.id = orders.shop_id").
		Where("shops.user_id=?", security.PayloadData.UserID).
		Find(&orders).Error
	return &orders, err
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id uint64, status string) error {

	order, err := r.GetOrderById(ctx, id)

	if err != nil {
		return err
	}

	order.Status = status
	return r.db.WithContext(ctx).Save(order).Error
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepositoryInterface {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if err := r.db.WithContext(ctx).
		Preload("OrderDetail").
		Create(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) GetOrderById(ctx context.Context, id uint64) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.WithContext(ctx).Preload("OrderDetail").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) DeleteOrder(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Order{}, id).Error
}
