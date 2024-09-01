package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func (s shopRepository) GetProducts(ctx context.Context, shopId uint64, userID uint64) (*[]entity.Product, error) {
	var product []entity.Product

	if err := s.db.WithContext(ctx).
		Model(&entity.Product{}).
		Joins("shops ON shops.id = products.shop_id").
		Where("products.shop_id = ? AND shops.user_id = ?", shopId, userID).
		Find(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

func (s shopRepository) GetOrders(ctx context.Context, shopID uint64, userID uint64) (*[]entity.Order, error) {
	var orders []entity.Order

	if err := s.db.WithContext(ctx).
		Model(&entity.Order{}).
		Joins("shops ON shops.id = orders.shop_id").
		Where("orders.shop_id=? AND shops.user_id=?", shopID, userID).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return &orders, nil

}

func (s shopRepository) List(ctx context.Context, userID uint64) (*[]entity.Shop, error) {
	var shop []entity.Shop

	if err := s.db.WithContext(ctx).
		Model(&entity.Shop{}).
		Find(&shop, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &shop, nil
}

func (s shopRepository) GetById(ctx context.Context, id uint64, userID uint64) (*entity.Shop, error) {
	var shop entity.Shop
	if err := s.db.WithContext(ctx).
		Model(&entity.Shop{}).
		First(&shop, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

func NewShopRepository(db *gorm.DB) repository.ShopRepositoryInterface {
	return &shopRepository{
		db: db,
	}
}
