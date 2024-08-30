package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
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
		Group("products.id").
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func NewProductRepository(db *gorm.DB) repository.ProductRepositoryInterface {
	return &productRepository{db: db}
}
