package service

import (
	"context"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/entity"
	repoPkg "github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"github.com/ecommerce-api/pkg/service"
	"gorm.io/gorm"
)

type shopService struct {
	shopRepo repoPkg.ShopRepositoryInterface
}

func (s shopService) List(ctx context.Context) (*[]entity.Shop, error) {
	return s.shopRepo.List(ctx, security.PayloadData.UserID)
}

func (s shopService) GetByID(ctx context.Context, shopID uint64) (*entity.Shop, error) {
	return s.shopRepo.GetById(ctx, shopID, security.PayloadData.UserID)
}

func (s shopService) Orders(ctx context.Context, shopId uint64) (*[]entity.Order, error) {
	return s.shopRepo.GetOrders(ctx, shopId, security.PayloadData.UserID)
}

func (s shopService) Products(ctx context.Context, shopID uint64) (*[]entity.Product, error) {
	return s.shopRepo.GetProducts(ctx, shopID, security.PayloadData.UserID)
}

func NewShopService(db *gorm.DB) service.ShopServiceInterface {
	shopRepo := repository.NewShopRepository(db)

	return &shopService{
		shopRepo: shopRepo,
	}
}
