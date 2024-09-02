package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func (p productRepository) ListWithStock(ctx context.Context) (*[]dto.ProductResponse, error) {

	var results []dto.ProductResponse

	err := p.db.WithContext(ctx).
		Model(&entity.Product{}).
		Select("products.*, sum(warehouse_inventories.quantity) as stock").
		Joins("LEFT JOIN warehouse_inventories ON products.id = warehouse_inventories.product_id").
		Joins("LEFT JOIN warehouses ON warehouses.id = warehouse_inventories.warehouse_id").
		Where("warehouses.user_id=?", security.PayloadData.UserID).
		Group("products.id").
		Find(&results).Error

	return &results, err
}

func (p productRepository) List(ctx context.Context) (*[]entity.Product, error) {
	var products []entity.Product
	err := p.db.WithContext(ctx).Model(&entity.Product{}).Find(&products).Error
	return &products, err
}

func (p productRepository) GetById(ctx context.Context, id uint64) (*entity.Product, error) {
	var product entity.Product
	if err := p.db.WithContext(ctx).Model(&product).
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p productRepository) GetByIdWithStock(ctx context.Context, id uint64) (*dto.ProductResponse, error) {
	var product dto.ProductResponse
	if err := p.db.WithContext(ctx).
		Model(&entity.Product{}).
		Select("products.*, sum(warehouse_inventories.quantity) as stock").
		Joins("LEFT JOIN warehouse_inventories ON products.id = warehouse_inventories.product_id").
		Joins("LEFT JOIN warehouses ON warehouses.id = warehouse_inventories.warehouse_id").
		Where("warehouses.user_id=?", security.PayloadData.UserID).
		Group("products.id").
		Debug().
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p productRepository) GetByShop(ctx context.Context, shopId uint64, userId uint64) (*[]dto.ProductResponse, error) {
	var results []dto.ProductResponse

	err := p.db.WithContext(ctx).
		Model(&entity.Product{}).
		Select("products.*, sum(warehouse_inventories.quantity) as stock").
		Joins("LEFT JOIN shops ON products.shop_id = shops.id").
		Joins("LEFT JOIN warehouse_inventories ON products.id = warehouse_inventories.product_id").
		Where("shops.id = ? AND shop.user_id", shopId, userId).
		Find(&results).Error

	return &results, err
}

func NewProductRepository(db *gorm.DB) repository.ProductRepositoryInterface {
	return &productRepository{db: db}
}
