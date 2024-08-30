package service

import (
	"context"
	repositoryPkg "github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/service"
	"gorm.io/gorm"
)

type productService struct {
	repo          repository.ProductRepositoryInterface
	inventoryRepo repository.WarehouseInventoryRepositoryInterface
}

func (s *productService) List(ctx context.Context) (*[]entity.Product, error) {
	return s.repo.List(ctx)
}

func (s *productService) ListWithStock(ctx context.Context) (*[]dto.ProductResponse, error) {
	return s.repo.ListWithStock(ctx)
}

func (s *productService) GetByID(ctx context.Context, id uint64) (*entity.Product, error) {
	return s.repo.GetById(ctx, id)
}

func (s *productService) GetByIDWithStock(ctx context.Context, id uint64) (*dto.ProductResponse, error) {
	return s.repo.GetByIdWithStock(ctx, id)
}

func (s *productService) GetAvailabilityStock(ctx context.Context, productID uint64, shopID uint64) (int, error) {
	return s.inventoryRepo.GetAvailableStock(ctx, productID, shopID)
}

func NewProductService(db *gorm.DB) service.ProductServiceInterface {
	productRepo := repositoryPkg.NewProductRepository(db)
	inventoryRepo := repositoryPkg.NewWarehouseInventoryRepository(db)
	return &productService{repo: productRepo, inventoryRepo: inventoryRepo}
}
